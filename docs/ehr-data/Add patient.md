# Add patient

## Person

```
{
	// CONTEXT AND GENERIC DATA

    "/context/health_care_facility|name": ">NAME OF THE CLINIC<",
    "/context/health_care_facility|identifier": ">ID OF THEs CLINIC<",
    "/context/start_time":"2018-04-09T10:18:17.352+00:00", // When added to waiting line?
    "/context/end_time":"2015-04-09T11:18:17.352+00:00", // When commited
    // if present, add administrator
    "/context/participation:0|function":"Administrator",
    "/context/participation:0|name":">NAME OF THE ADMINISTRATOR<",
    "/context/participation:0|identifier":">USER ID OF THE ADMINSTRATOR<",
    "/context/participation:0|mode":"face-to-face communication::openehr::216",
    // if present, add nurse
    "/context/participation:1|function":"Nurse",
    "/context/participation:1|name":">NAME OF THE NURSE<",
    "/context/participation:1|identifier":">USER ID OF THE NURSE<",
    "/context/participation:1|mode":"face-to-face communication::openehr::216",
    "/composer|identifier":">DOCTOR'S (USER) ID<",
    "/composer|name":">NAME OF THE DOCTOR<",
    "/category":"openehr::431|persistent|",
    "/territory":">ISO CODE OF CLINIC'S COUNTRY<",
    "/language":"en",

    // PATIENT DATA

    // First / Given name (DV_TEXT)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0002]": ">GIVEN NAME<",

    // Last name (DV_TEXT)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0003]": ">FAMILY NAME<",

    // Preferred name (DV_BOOLEAN)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0008]": "true",

    // Date of birth (DV_DATE)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0010]": "true",

	// Country of Birth (CODED)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]": "countries::SI|SLOVENIA|",

	// Gender (CODED; at0310:"Male"; at0311:"Female"; at0312:"Intersex or indeterminate"; at0313:"Not declared/inadequately described")
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]": "local::at0310|Male|",

	// Marital Status (CODED; Childer of SNOMED 365581002)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]": "SNOMED::87915002|Married|",

	// Example Identity (list, ex: Syrian ID)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:0|id": ">ID NUMBER<",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:0|type": ">ID TYPE<",

	/* ADDRESS */
	// Type of address
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0033]": "local::at0463|Temporary Accomodation|",
	// Country
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at0009]": "countries::SI|SLOVENIA|",
    // Camp (stored under "Address site name")
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at00014]": ">ID OF THE CAMP<",
    // Tent (stored under "Building/Complex sub-unit number")
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at00013]": ">TENT NUMBER<",

	/* Phone number */
	// Type of address (phone)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/name[at0014]": "local::at0022|Mobile|",
	"/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/details[at0001]/items[at0007]": ">MOBILE PHONE NUMBER<",

	/* Email address */
	// Type of address (email)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/name[at0014]": "local::at0024|Email|",
	"/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/details[at0001]/items[at0007]": ">EMAIL ADDRESS<",

	/* Whatsapp */
	// Type of address (whatsapp)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:3/name[at0013]": "whatsapp",
	"/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/details[at0001]/items[at0007]": ">WHATSAPP PHONE NUMBER<",

	// Relationships
	// type of relationship
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0/details[at0041]": ">relationships::stranger|Thing|<",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0|namespace": "local-together", (local-together, local)
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0|type": "PERSON",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0|id": ">PERSON-UUID<"
}
```

## Patient Info

```
{
	// CONTEXT AND GENERIC DATA

    "/context/health_care_facility|name": ">NAME OF THE CLINIC<",
    "/context/health_care_facility|identifier": ">ID OF THEs CLINIC<",
    "/context/start_time":"2018-04-09T10:18:17.352+00:00", // When added to waiting line?
    "/context/end_time":"2015-04-09T11:18:17.352+00:00", // When commited
    // if present, add administrator
    "/context/participation:0|function":"Administrator",
    "/context/participation:0|name":">NAME OF THE ADMINISTRATOR<",
    "/context/participation:0|identifier":">USER ID OF THE ADMINSTRATOR<",
    "/context/participation:0|mode":"face-to-face communication::openehr::216",
    // if present, add nurse
    "/context/participation:0|function":"Nurse",
    "/context/participation:0|name":">NAME OF THE NURSE<",
    "/context/participation:0|identifier":">USER ID OF THE NURSE<",
    "/context/participation:0|mode":"face-to-face communication::openehr::216",
    "/composer|identifier":">DOCTOR'S (USER) ID<",
    "/composer|name":">NAME OF THE DOCTOR<",
    "/category":"openehr::431|persistent|",
    "/territory":">ISO CODE OF CLINIC'S COUNTRY<",
    "/language":"en",

    // MEDICAL HISTORY

    /* Chronic diseases (array) */
    // Name
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:0/items[at0018]": ">NAME OF THE DISEASE<",
    // Date
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:0/items[at0017]": ">DATE (YYYY-MM-DD)<",
    // Comment
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:0/items[at0016]": ">COMMENT<",

    /* Immunisations (array) */
    // Name
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:0/items[at0019]": ">NAME<",
    // Date
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:0/items[at0021]": ">DATE (YYYY-MM-DD)<",
    // Comment
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:0/items[at0027]": ">COMMENT<",

    /* Allergies (array) */
    // Name
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:0/items[at0010]": ">NAME OF THE ALLERGY<",
    // Critical (BOOL)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at00009]:0/items[at0012]": ">true|false<",
    // Comment
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:0/items[at0013]": ">COMMENT<",

    /* Injuries or handicaps (array) */
    // Name
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:0/items[at0023]": ">NAME<",
    // Date
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:0/items[at0024]": ">DATE (YYYY-MM-DD)<",
    // Comment
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:0/items[at0025]": ">COMMENT<",

    /* Surgeries (array) */
    // Name
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:0/items[at0028]": ">NAME<",
    // Date
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:0/items[at0029]": ">DATE (YYYY-MM-DD)<",
    // Comment
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:0/items[at0030]": ">COMMENT<",

    /* Medications (array) */
    // Name
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:0/items[at0032]": ">NAME<",
    // Date
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:0/items[at0033]": ">DATE (YYYY-MM-DD)<",
    // Comment
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:0/items[at0034]": ">COMMENT<",

    /* Kids and adults questionaire*/

    // Are you smoking (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0039]": ">true|false<",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0038]": ">comment<",

    // Taking drugs (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0051]/items[at0040]": ">true|false<",

    // Resources for basic hygiene (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0052]/items[at0041]": ">true|false<",

    // Access to clean water (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0053]/items[at0042]": ">true|false<",

    // Sufficient food supply (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0054]/items[at0043]": ">true|false<",

    // Good appetite (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0055]/items[at0044]": ">true|false<",

    // Accomodations have heating (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0056]/items[at0045]": ">true|false<",


    // Accomodations have electricity (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0057]/items[at0046]": ">true|false<",

	/* Additional patient info */

    // Number of kids (NUMBER)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0047]": >NUMBER<,

    // Nationality (CODED)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0048]": ">countries::SI|Slovenia|<",

    // Country of origin (CODED)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0049]": ">countries::SI|Slovenia|<",

    // Education (CODED)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0050]": ">education::secondary|Secondary education|<",

    // Occupation
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0058]": ">OCCUPATION<",

    // Date of leaving home country (DATE)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0059]": ">DATE(YYYY-MM-DD)<",

    // Date of arriving to camp (DATE)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0061]": ">DATE(YYYY-MM-DD)<",

    // Transit countries
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0060]": ">TRANSIT COUNTRIES<",

    /* BABY VACCINE INFORMATION */

	// On schedule at home (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0065]": ">true|false<",

	// Has Immunization documents (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0066]": ">true|false<",

	// Tested for tuberculosis (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0067]": ">true|false<",

	// Were tuberculosis tests positivee (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0068]": ">true|false<",

	// Any additional tests done (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0069]": ">true|false<",

	// Investigation details
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0070]": ">DETAILS<",

	// Any reaction to vaccines (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0071]": ">true|false<",

	// Details of vaccine reactions
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0072]": ">DETAILS<",

    /* Basic baby screening */

	// Delivery type (CODED)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0075]": ">deliveryTypes::something|SOMETHING|<",

    // Premature (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0076]": ">true|false<",

    // Weeks at birth (INTEGER)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0077]": >WEEKS<,

    // Weight at birth (QUANTITY)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0078]": ">WEIGHT<,gm",

    // Height at birth (QUANTITY)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0079]": ">HEIGHT<,cm",

    // Breastfeeding (BOOL)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0087]/items[at0081]": ">true|false<",

    // Breastfeeding for how long (INTEGER)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0087]/items[at0082]": >NUM<,

    // What does baby eat or drink (CODED)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0083]": ">babyFood::something|Something|<",

    // How many diapers does child wet (INTEGER)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0084]": >NUM<,

    // How many times does child have bowl movement (CODED)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0085]": ">bowelMovementFrequency::2perweek|Twice per week<",

    // Describe bowl movement (INTEGER)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0086]": ">DESCRIPTION<",

    // Satisfied with sleep (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0088]/items[at0089]": ">true|false<",

    // Comment about the sleep
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0088]/items[at0090]": ">COMMENT<",

    // Do you take / give vit. D (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0091]": ">true|false<",

    // Baby sleeps on her back (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0092]": ">true|false<",

    // Does anyone smoke (BOOLEAN)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0093]": ">true|false<",

    // Number of smokers (INTEGER)
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0094]": >NUM<,

    // How does child get around
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0095]": ">TXT<",

    // How does child communicate
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0096]": ">TXT<"
}
```



















