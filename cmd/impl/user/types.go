package user

// RoleLookup maps role ids to friendly names
// https://www.site24x7.com/help/api/#user_constants
var RoleLookup = map[int]string{
	0: "No Access",
	1: "Super Administrator",
	2: "Administrator",
	3: "Operator",
	4: "Billing Contact",
	5: "Spokesperson",
	6: "Hosting Provider",
	7: "Read Only",
}

// CloudspendRoles maps role ids to friendly names
// https://www.site24x7.com/help/api/#user_constants
var CloudspendRoles = map[int]string{
	11: "Cost Administrator",
	12: "Cost User",
}

// StatusIQRoles maps role ids to friendly names
// https://www.site24x7.com/help/api/#alerting_constants
var StatusIQRoles = map[int]string{
	21: "StatusIQ Super Administrator",
	22: "StatusIQ Administrator",
	23: "StatusIQ SpokesPerson",
	24: "StatusIQ Billing Contact",
	25: "StatusIQ Read Only",
}

// NotificationMethods are communication channels through which alerts can be sent
// https://www.site24x7.com/help/api/#alerting_constants
var NotificationMethods = map[int]string{
	1: "Email",
	2: "SMS",
	3: "Voice Call",
	4: "IM",
	5: "Twitter",
}

// EmailFormats are the possible formats that can be used for notification emails
// https://www.site24x7.com/help/api/#alerting_constants
var EmailFormats = map[int]string{
	0: "TEXT",
	1: "HTML",
}

// JobTitles are the supported options for user job title
// https://www.site24x7.com/help/api/#job_title
var JobTitles = map[int]string{
	1: "IT Engineer",
	2: "Cloud Engineer",
	3: "DevOps Engineer",
	4: "Webmaster",
	5: "CEO/CTO",
	6: "Internal IT",
	7: "Others",
}

// ResourceTypes are the supported resource/selection types
// https://www.site24x7.com/help/api/#resource_type_constants
var ResourceTypes = map[int]string{
	0: "All Monitors",
	1: "Monitor Group",
	2: "Monitor",
	3: "Tags",
	4: "Monitor Type",
}
