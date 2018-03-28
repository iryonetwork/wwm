const requiredFields = [
    "firstName",
    "lastName",
    "dateOfBirth",
    "gender",
    "maritalStatus",
    "numberOfKids",
    "nationality",
    "countryOfOrigin",
    "country",
    "camp",
    "tent"
]

const validate = values => {
    const errors = {}
    /*
    if (!values.firstName) {
        errors.firstName = "Required"
    }
    if (!values.lastName) {
        errors.lastName = "Required"
    }
    */

    requiredFields.forEach(field => {
        if (!values[field] || values[field].trim() === "") {
            errors[field] = "Required"
        }
    })

    return errors
}

export default validate
