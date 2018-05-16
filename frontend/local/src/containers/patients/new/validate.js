const requiredFields = ["firstName", "lastName", "dateOfBirth", "gender", "maritalStatus", "numberOfKids", "nationality", "countryOfOrigin", "country", "camp"]

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
