import { load as loadCode } from "shared/modules/codes"

export default dispatch => {
    return Promise.all([dispatch(loadCode("countries")), dispatch(loadCode("gender")), dispatch(loadCode("maritalStatus"))]).then(
        ([countries, genders, maritalStatuses]) => [
            // Name
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0002]",
                formPath: "firstName"
            },
            // Last name
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0003]",
                formPath: "lastName"
            },
            // Preferred Name
            {
                type: "fixedValue",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0008]",
                value: "true"
            },
            // Date of Birth
            {
                type: "value",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0010]",
                formPath: "dateOfBirth"
            },
            // Country of Birth
            {
                type: "code",
                codes: countries,
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]",
                formPath: "countryOfOrigin"
            },
            // Gender
            {
                type: "code",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0017]",
                formPath: "gender",
                codes: genders
            },
            // Marital Status
            {
                type: "code",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0033]",
                formPath: "maritalStatus",
                codes: maritalStatuses
            },
            /* ADDRESS */
            // Type of address
            {
                type: "fixedValue",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0033]",
                value: "local::at0463|Temporary Accommodation|"
            },
            // Country
            {
                type: "code",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at0009]",
                formPath: "country",
                codes: countries
            },
            // Camp (stored under "Address site name")
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at00014]",
                formPath: "camp"
            },
            // Tent (stored under "Building/Complex sub-unit number")
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at00013]",
                formPath: "tent"
            },

            /* Phone number */
            // Type of address (phone)
            {
                type: "fixedValue",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/name[at0014]",
                value: "local::at0022|Mobile|"
            },
            // Phone number
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/details[at0001]/items[at0007]",
                formPath: "phone"
            },

            /* Email address */
            // Type of address (email)
            {
                type: "fixedValue",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/name[at0014]",
                value: "local::at0024|Email|"
            },
            // Email address
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/details[at0001]/items[at0007]",
                formPath: "email"
            },

            /* Whatsapp */
            // Type of address (whatsapp)
            {
                type: "fixedValue",
                ehrPath: "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:3/name[at0013]",
                value: "whatsapp"
            },
            // Address
            {
                type: "value",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:3/details[at0001]/items[at0007]",
                formPath: "whatsapp"
            },
            // Documents
            {
                type: "array",
                ehrPath:
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]",
                formPath: "documents",
                items: [
                    {
                        type: "value",
                        ehrPath: "|id",
                        formPath: "number"
                    },
                    {
                        type: "value",
                        ehrPath: "|type",
                        formPath: "type"
                    }
                ]
            }
        ]
    )
}
