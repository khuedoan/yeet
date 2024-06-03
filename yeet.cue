// Kitchen sink example for all options

// By default if no yeet.cue is specified but webhook is set up, it will only build the app.
// It will automatically check if that commit is already built, if not it will build it
// It will tag the image as git reference (e.g. master) AND commit hash,
//
// For example when you push to master the following image tags are pushed:
// - khuedoan/example-service:master
// - khuedoan/example-service:75d71cc66d6a7f8815fa30978089c862046edace

// This is a regex for Git reference
// Usually it should trigger on master, but it doesn't have to be
"master": {
    // Because dev environments tracks a static branch name,
    // we should include a hash to update pod image version.
    // In GitOps config files (e.g. Helm values, Timoni values), this will be set as:
    // image:
    //   repository: khuedoan/example-service
    //   tag: master@sha256:5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03
    includeHash: true // default false

    // Where do you wanna deploy?
    stages: [
        // When pushing to master there's only a single stage to deploy to AG1
        {
            name: "dev"
            groups: [
                // In each stage, there's multiple groups of targets to deploy to
                // Each groups must complete before continuing to the next one
                // Targets inside a group will be deployed in parallel
                //
                // Usually a target is a cluster, it will recursively find the targets based on the group
                "dev/main",
                "dev/dr"
            ]
        }
    ]
}

// Release strategy for a git tag, usually on a good commit in master
"v*.*.*": {
    stages: [
        // Now we have multiple stages after releasing the app
        // First deploy to staging
        {
            name: "staging"
            groups: [
                "stg"
            ]
            // Wait 1 day before continuing to the next stage
            wait: "1d"
        },
        // We wanna deploy to low risk production groups first
        {
            name: "production-low"
            groups: [
                "vn1",
                "in1",
                "nl1",
                "eu1"
            ]
            // Wait a bit then continue to the next stage
            // Usually this must be enough for QA engineers to run tests
            wait: "2h"
        },
        // Same as the previous one
        {
            name: "production-medium"
            groups: [
                "uk1",
                "la1"
            ]
            // And wait more for smoke tests and let it bake, and come back tomorrow
            wait: "1d"
        },
        // The next day, continue to higher risk production groups
        {
            name: "production-high"
            groups: [
                // A groups can be as specific as it needs, down to the cluster level...
                "us1/main/west/cluster-a",
                "us1/main/west/cluster-b",
                // ...or just the deployment name, it will recursively find member targets
                "usf"
            ]
        }
    ]
}
