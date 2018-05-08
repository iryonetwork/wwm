export const getName = user => {
    if (user.personalData) {
        var name = ""
        if (user.personalData.firstName !== undefined && user.personalData.firstName !== "") {
            name += user.personalData.firstName
        }
        if (user.personalData.middleName !== undefined && user.personalData.middleName !== "") {
            name = name + " " + user.personalData.middleName
        }
        if (user.personalData.lastName !== undefined && user.personalData.lastName !== "") {
            name = name + " " + user.personalData.lastName
        }
        return name
    }

    return "Unknown"
}
