const requiredFields = [
    "firstName",
    "lastName",
    "dateOfBirth",
    "gender",
    "numberOfKids",
    "nationality",
    "countryOfOrigin",
    "country",
    "region",
]

const validate = values => {
    const errors = {}
    requiredFields.forEach(field => {
        if (!values[field] || values[field].trim() === "") {
            errors[field] = "Required"
        }
    })

    return errors
}

export default validate
