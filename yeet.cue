// Kitchen sink example for all options

// By default if no yeet.yaml is specified but webhook is set up,
// it will use the following default config:
// push: {
//     stages: [{
//         name: "dev"
//         groups: [
//             "dev"
//         ]
//     }]
// }

// This is an event
// Usually it should trigger on master, but it doesn't have to be
push: {
    // Wanna build and push using the event metadata?
    // It will tag the image as branch name (master) AND commit hash
    build: true  // default

    // After build where do you wanna deploy?
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
tag: {
    // A push to master already built
    // Just tag the commit hash to a new name, e.g. v1.2.3
    build: false

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
                // ...or just the deployment name
                "usf"
            ]
        }
    ]
}
