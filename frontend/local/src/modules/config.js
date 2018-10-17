import { LENGTH_UNIT, WEIGHT_UNIT, TEMPERATURE_UNIT, BLOOD_PRESSURE_UNIT } from "shared/unitConversion/units"

export const READ_ONLY_KEYS = "readOnlyKeys"
export const LOCALE = "locale"
export const BASE_URL = "baseUrl"
export const API_URL = "apiUrl"
export const CLINIC_ID = "clinicId"
export const LOCATION_ID = "locationId"
export const BABY_MAX_AGE = "babyMaxAge"
export const CHILD_MAX_AGE = "childMaxAge"
export const DEFAULT_WAITLIST_ID = "waitlistId"
export { LENGTH_UNIT, WEIGHT_UNIT, TEMPERATURE_UNIT, BLOOD_PRESSURE_UNIT }

export const localeOptions = [{ value: "en", label: "English" }]
export const lengthUnitOptions = [{ value: "cm", label: "Centimeters" }, { value: "ft_in", label: "Feet and inches" }]
export const weightUnitOptions = [{ value: "kg_g", label: "Kilograms and grams" }, { value: "lb_oz", label: "Pounds and ounces" }]
export const temperatureUnitOptions = [{ value: "°C", label: "Celsius" }, { value: "°F", label: "Fahrenheit" }]
export const bloodPressureUnitOptions = [{ value: "mm[Hg]", label: "mm[Hg]" }, { value: "cm[Hg]", label: "cm[Hg]" }]
