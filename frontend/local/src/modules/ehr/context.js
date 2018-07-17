export default dispatch =>
    Promise.resolve([
        {
            type: "value",
            ehrPath: "/context/health_care_facility|name",
            formPath: "context.clinic.name"
        },
        {
            type: "value",
            ehrPath: "/context/health_care_facility|identifier",
            formPath: "context.clinic.id"
        },
        {
            type: "value",
            ehrPath: "/territory",
            formPath: "context.territory"
        },
        {
            type: "value",
            ehrPath: "/language",
            formPath: "context.language"
        },
        {
            type: "dateTime",
            ehrPath: "/context/start_time",
            formPath: "context.startTime"
        },
        {
            type: "dateTime",
            ehrPath: "/context/end_time",
            formPath: "context.endTime"
        },
        {
            type: "value",
            ehrPath: "/composer|identifier",
            formPath: "context.author.id"
        },
        {
            type: "value",
            ehrPath: "/composer|name",
            formPath: "context.author.name"
        },
        {
            type: "value",
            ehrPath: "/category",
            formPath: "context.category"
        }
    ])
