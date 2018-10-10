export const required = (value, values, props) => {
    return value && value !== "" ? undefined : "Required"
}
export const expectedRange = (min, max) => value => {
    return value < min || value > max ? "Entered value is out of expected range" : undefined
}

export const babyHeightExpectedRange = expectedRange(25, 80)
export const heightExpectedRange = expectedRange(0, 250)
export const babyWeightExpectedRange = expectedRange(1000, 6000)
export const weightExpectedRange = expectedRange(0, 250)
export const temperatureExpectedRange = expectedRange(30, 50)
export const systolicBloodPressureExpectedRange = expectedRange(60, 200)
export const diastolicBloodPressureExpectedRange = expectedRange(30, 110)
export const heartRateExpectedRange = expectedRange(30, 250)
export const oxygenSaturationxpectedRange = expectedRange(70, 100)
