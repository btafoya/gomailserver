package webdav

import (
	"encoding/xml"
	"time"
)

// PropFind represents a PROPFIND request
type PropFind struct {
	XMLName  xml.Name  `xml:"DAV: propfind"`
	AllProp  *struct{} `xml:"allprop"`
	PropName *struct{} `xml:"propname"`
	Prop     *Prop     `xml:"prop"`
}

// Prop represents requested properties
type Prop struct {
	ResourceType         *struct{} `xml:"DAV: resourcetype"`
	DisplayName          *struct{} `xml:"DAV: displayname"`
	GetContentType       *struct{} `xml:"DAV: getcontenttype"`
	GetETag              *struct{} `xml:"DAV: getetag"`
	GetLastModified      *struct{} `xml:"DAV: getlastmodified"`
	GetContentLength     *struct{} `xml:"DAV: getcontentlength"`
	CreationDate         *struct{} `xml:"DAV: creationdate"`
	CurrentUserPrincipal *struct{} `xml:"DAV: current-user-principal"`
	// CalDAV specific
	CalendarData       *CalendarDataRequest `xml:"urn:ietf:params:xml:ns:caldav calendar-data"`
	CalendarHomeSet    *struct{}            `xml:"urn:ietf:params:xml:ns:caldav calendar-home-set"`
	CalendarDescription *struct{}           `xml:"urn:ietf:params:xml:ns:caldav calendar-description"`
	CalendarColor      *struct{}            `xml:"http://apple.com/ns/ical/ calendar-color"`
	CalendarOrder      *struct{}            `xml:"http://apple.com/ns/ical/ calendar-order"`
	SupportedCalendarComponentSet *struct{} `xml:"urn:ietf:params:xml:ns:caldav supported-calendar-component-set"`
	// CardDAV specific
	AddressData        *AddressDataRequest `xml:"urn:ietf:params:xml:ns:carddav address-data"`
	AddressbookHomeSet *struct{}           `xml:"urn:ietf:params:xml:ns:carddav addressbook-home-set"`
	AddressbookDescription *struct{}       `xml:"urn:ietf:params:xml:ns:carddav addressbook-description"`
	SupportedAddressData *struct{}         `xml:"urn:ietf:params:xml:ns:carddav supported-address-data"`
}

// CalendarDataRequest represents a calendar-data property request
type CalendarDataRequest struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:caldav calendar-data"`
	// Optional filters and component selection
}

// AddressDataRequest represents an address-data property request
type AddressDataRequest struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:carddav address-data"`
	// Optional filters and property selection
}

// MultiStatus represents a WebDAV multistatus response
type MultiStatus struct {
	XMLName   xml.Name   `xml:"DAV: multistatus"`
	Responses []Response `xml:"response"`
	SyncToken string     `xml:"sync-token,omitempty"`
}

// Response represents a single resource response
type Response struct {
	Href      string     `xml:"href"`
	PropStats []PropStat `xml:"propstat"`
	Status    string     `xml:"status,omitempty"`
	Error     *Error     `xml:"error,omitempty"`
}

// PropStat represents property status
type PropStat struct {
	Prop   PropValue `xml:"prop"`
	Status string    `xml:"status"`
}

// PropValue represents property values in a response
type PropValue struct {
	ResourceType         *ResourceType `xml:"DAV: resourcetype,omitempty"`
	DisplayName          *string       `xml:"DAV: displayname,omitempty"`
	GetContentType       *string       `xml:"DAV: getcontenttype,omitempty"`
	GetETag              *string       `xml:"DAV: getetag,omitempty"`
	GetLastModified      *string       `xml:"DAV: getlastmodified,omitempty"`
	GetContentLength     *int64        `xml:"DAV: getcontentlength,omitempty"`
	CreationDate         *string       `xml:"DAV: creationdate,omitempty"`
	CurrentUserPrincipal *Href         `xml:"DAV: current-user-principal,omitempty"`
	// CalDAV specific
	CalendarData       *string           `xml:"urn:ietf:params:xml:ns:caldav calendar-data,omitempty"`
	CalendarHomeSet    *Href             `xml:"urn:ietf:params:xml:ns:caldav calendar-home-set,omitempty"`
	CalendarDescription *string          `xml:"urn:ietf:params:xml:ns:caldav calendar-description,omitempty"`
	CalendarColor      *string           `xml:"http://apple.com/ns/ical/ calendar-color,omitempty"`
	CalendarOrder      *int              `xml:"http://apple.com/ns/ical/ calendar-order,omitempty"`
	SupportedCalendarComponentSet *SupportedCalendarComponentSet `xml:"urn:ietf:params:xml:ns:caldav supported-calendar-component-set,omitempty"`
	// CardDAV specific
	AddressData        *string `xml:"urn:ietf:params:xml:ns:carddav address-data,omitempty"`
	AddressbookHomeSet *Href   `xml:"urn:ietf:params:xml:ns:carddav addressbook-home-set,omitempty"`
	AddressbookDescription *string `xml:"urn:ietf:params:xml:ns:carddav addressbook-description,omitempty"`
	SupportedAddressData *SupportedAddressData `xml:"urn:ietf:params:xml:ns:carddav supported-address-data,omitempty"`
}

// ResourceType represents the type of a WebDAV resource
type ResourceType struct {
	Collection          *struct{} `xml:"DAV: collection,omitempty"`
	Calendar            *struct{} `xml:"urn:ietf:params:xml:ns:caldav calendar,omitempty"`
	Addressbook         *struct{} `xml:"urn:ietf:params:xml:ns:carddav addressbook,omitempty"`
	Principal           *struct{} `xml:"DAV: principal,omitempty"`
}

// Href represents a reference to a resource
type Href struct {
	Href string `xml:"href"`
}

// SupportedCalendarComponentSet represents supported calendar component types
type SupportedCalendarComponentSet struct {
	Components []CalendarComponent `xml:"comp"`
}

// CalendarComponent represents a calendar component type
type CalendarComponent struct {
	Name string `xml:"name,attr"`
}

// SupportedAddressData represents supported address data formats
type SupportedAddressData struct {
	AddressDataTypes []AddressDataType `xml:"address-data-type"`
}

// AddressDataType represents a supported vCard format
type AddressDataType struct {
	ContentType string `xml:"content-type,attr"`
	Version     string `xml:"version,attr"`
}

// Error represents a WebDAV error
type Error struct {
	XMLName xml.Name `xml:"DAV: error"`
	// Error conditions
}

// FormatHTTPDate formats a time for HTTP headers (RFC 1123)
func FormatHTTPDate(t time.Time) string {
	return t.UTC().Format(time.RFC1123)
}

// FormatISO8601 formats a time in ISO 8601 format
func FormatISO8601(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
