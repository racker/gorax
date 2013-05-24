// vim: ts=8 sw=8 noet ai

package servers

/*** Common Sub-elements ***/

// Link is used for JSON (un)marshalling.
// It provides RESTful links to a resource.
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

// FileConfig structures represent a blob of data which must appear at a
// a specific location in a server's filesystem.  The file contents are
// base-64 encoded.
type FileConfig struct {
	Path     string `json:"path"`
	Contents string `json:"contents"`
}

// NetworkConfig structures represent an affinity between a server and a
// specific, uniquely identified network.  Networks are identified through
// universally unique IDs.
type NetworkConfig struct {
	Uuid string `json:"uuid"`
}

/*** Images ***/

// Image is used for JSON (un)marshalling.
// It provides a description of an OS image.
//
// The Id field contains the image's unique identifier.
// For example, this identifier will be useful for specifying which operating system to install on a new server instance.
//
// The MinDisk and MinRam fields specify the minimum resources a server must provide to be able to install the image.
//
// The Name field provides a human-readable moniker for the OS image.
//
// The Progress and Status fields indicate image-creation status.
// Any usable image will have 100% progress.
//
// The Updated field indicates the last time this image was changed.
//
// OsDcfDiskConfig indicates the server's boot volume configuration.
// Valid values are:
//     AUTO
//     ----
//     The server is built with a single partition the size of the target flavor disk.
//     The file system is automatically adjusted to fit the entire partition.
//     This keeps things simple and automated.
//     AUTO is valid only for images and servers with a single partition that use the EXT3 file system.
//     This is the default setting for applicable Rackspace base images.
//
//     MANUAL
//     ------
//     The server is built using whatever partition scheme and file system is in the source image.
//     If the target flavor disk is larger,
//     the remaining disk space is left unpartitioned.
//     This enables images to have non-EXT3 file systems, multiple partitions, and so on,
//     and enables you to manage the disk configuration.
//
type Image struct {
	OsDcfDiskConfig string `json:"OS-DCF:diskConfig"`
	Created         string `json:"created"`
	Id              string `json:"id"`
	Links           []Link `json:"links"`
	MinDisk         int    `json:"minDisk"`
	MinRam          int    `json:"minRam"`
	Name            string `json:"name"`
	Progress        int    `json:"progress"`
	Status          string `json:"status"`
	Updated         string `json:"updated"`
}

// ImageLink provides a reference to a image by either ID or by direct URL.
// Some services use just the ID, others use just the URL.
// This structure provides a common means of expressing both in a single field.
type ImageLink struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

/*** Flavors ***/

// Flavor records represent (virtual) hardware configurations for server resources in a region.
//
// The Id field contains the flavor's unique identifier.
// For example, this identifier will be useful when specifying which hardware configuration to use for a new server instance.
//
// The Disk and Ram fields provide a measure of storage space offered by the flavor, in GB and MB, respectively.
//
// The Name field provides a human-readable moniker for the flavor.
//
// Swap indicates how much space is reserved for swap.
// If not provided, this field will be set to 0.
//
// VCpus indicates how many (virtual) CPUs are available for this flavor.
type Flavor struct {
	OsFlvDisabled bool    `json:"OS-FLV-DISABLED:disabled"`
	Disk          int     `json:"disk"`
	Id            string  `json:"id"`
	Links         []Link  `json:"links"`
	Name          string  `json:"name"`
	Ram           int     `json:"ram"`
	RxTxFactor    float64 `json:"rxtx_factor"`
	Swap          int     `json:"swap"`
	VCpus         int     `json:"vcpus"`
}

// FlavorLink provides a reference to a flavor by either ID or by direct URL.
// Some services use just the ID, others use just the URL.
// This structure provides a common means of expressing both in a single field.
type FlavorLink struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

/*** Servers ***/

// A VersionedAddress denotes either an IPv4 or IPv6 (depending on version indicated)
// address.
type VersionedAddress struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
}

// An AddressSet provides a set of public and private IP addresses for a resource.
// Each address has a version to identify if IPv4 or IPv6.
type AddressSet struct {
	Public  []VersionedAddress `json:"public"`
	Private []VersionedAddress `json:"private"`
}

// Server records represent (virtual) hardware instances (not configurations) accessible by the user.
//
// The AccessIPv4 / AccessIPv6 fields provides IP addresses for the server in the IPv4 or IPv6 format, respectively.
//
// Addresses provides addresses for any attached isolated networks
// and Rackspace public and private networks.
// The version field indicates whether the IP address is version 4 or 6.
//
// Created tells when the server entity was created.
//
// The Flavor field includes the flavor ID and flavor links.
//
// The compute provisioning algorithm has an anti-affinity property that
// attempts to spread customer VMs across hosts.
// Under certain situations,
// VMs from the same customer might be placed on the same host.
// The HostId field represents the host your server runs on and
// can be used to determine this scenario if it is relevant to your application.
//
// HostId is unique per account and is not globally unique.
// 
// Id provides the server's unique identifier.
// This field must be treated opaquely.
//
// Image indicates which image is installed on the server.
//
// Links provides one or more means of accessing the server.
//
// Metadata provides a small key-value store for application-specific information.
//
// Name provides a human-readable name for the server.
//
// Progress indicates how far along it is towards being provisioned.
// 100 represents complete, while 0 represents just beginning.
//
// Status provides an indication of what the server's doing at the moment.
// A server will be in ACTIVE state if it's ready for use.
//
// OsDcfDiskConfig indicates the server's boot volume configuration.
// Valid values are:
//     AUTO
//     ----
//     The server is built with a single partition the size of the target flavor disk.
//     The file system is automatically adjusted to fit the entire partition.
//     This keeps things simple and automated.
//     AUTO is valid only for images and servers with a single partition that use the EXT3 file system.
//     This is the default setting for applicable Rackspace base images.
//
//     MANUAL
//     ------
//     The server is built using whatever partition scheme and file system is in the source image.
//     If the target flavor disk is larger,
//     the remaining disk space is left unpartitioned.
//     This enables images to have non-EXT3 file systems, multiple partitions, and so on,
//     and enables you to manage the disk configuration.
//
// RaxBandwidth provides measures of the server's inbound and outbound bandwidth per interface.
//
// OsExtStsPowerState provides an indication of the server's power.
// This field appears to be a set of flag bits:
//
//           ... 4  3   2   1   0
//         +--//--+---+---+---+---+
//         | .... | 0 | S | 0 | I |
//         +--//--+---+---+---+---+
//                      |       |
//                      |       +---  0=Instance is down.
//                      |             1=Instance is up.
//                      |
//                      +-----------  0=Server is switched ON.
//                                    1=Server is switched OFF.
//                                    (note reverse logic.)
//
// Unused bits should be ignored when read, and written as 0 for future compatibility.
//
// OsExtStsTaskState and OsExtStsVmState work together
// to provide visibility in the provisioning process for the instance.
// Consult Rackspace documentation at
// http://docs.rackspace.com/servers/api/v2/cs-devguide/content/ch_extensions.html#ext_status
// for more details.  It's too lengthy to include here.
type Server struct {
	AccessIPv4         string         `json:"accessIPv4"`
	AccessIPv6         string         `json:"accessIPv6"`
	Addresses          AddressSet     `json:"addresses"`
	Created            string         `json:"created"`
	Flavor             FlavorLink     `json:"flavor"`
	HostId             string         `json:"hostId"`
	Id                 string         `json:"id"`
	Image              ImageLink      `json:"image"`
	Links              []Link         `json:"links"`
	Metadata           interface{}    `json:"metadata"`
	Name               string         `json:"name"`
	Progress           int            `json:"progress"`
	Status             string         `json:"status"`
	TenantId           string         `json:"tenant_id"`
	Updated            string         `json:"updated"`
	UserId             string         `json:"user_id"`
	OsDcfDiskConfig    string         `json:"OS-DCF:diskConfig"`
	RaxBandwidth       []RaxBandwidth `json:"rax-bandwidth:bandwidth"`
	OsExtStsPowerState int            `json:"OS-EXT-STS:power_state"`
	OsExtStsTaskState  string         `json:"OS-EXT-STS:task_state"`
	OsExtStsVmState    string         `json:"OS-EXT-STS:vm_state"`
}

// NewServer structures are used for both requests and responses.
// The fields discussed below are relevent for server-creation purposes.
//
// The Name field contains the desired name of the server.
// Note that (at present) Rackspace permits more than one server with the same name;
// however, software should not depend on this.
// Not only will Rackspace support thank you, so will your own devops engineers.
// A name is required.
//
// The ImageRef field contains the ID of the desired software image to place on the server.
// This ID must be found in the image slice returned by the Images() function.
// This field is required.
//
// The FlavorRef field contains the ID of the server configuration desired for deployment.
// This ID must be found in the flavor slice returned by the Flavors() function.
// This field is required.
//
// For OsDcfDiskConfig, refer to the Image or Server structure documentation.
// This field defaults to "AUTO" if not explicitly provided.
//
// Metadata contains a small key/value association of arbitrary data.
// Neither Rackspace nor OpenStack places significance on this field in any way.
// This field defaults to an empty map if not provided.
//
// Personality specifies the contents of certain files in the server's filesystem.
// The files and their contents are mapped through a slice of FileConfig structures.
// If not provided, all filesystem entities retain their image-specific configuration.
//
// Networks specifies an affinity for the server's various networks and interfaces.
// Networks are identified through UUIDs; see NetworkConfig structure documentation for more details.
// If not provided, network affinity is determined automatically.
//
// The AdminPass field may be used to provide a root- or administrator-password
// during the server provisioning process.
// If not provided, a random password will be automatically generated and returned in this field.
//
// The following fields are intended to be used to communicate certain results about the server being provisioned.
// When attempting to create a new server, these fields MUST not be provided.
// They'll be filled in by the response received from the Rackspace APIs.
//
// The Id field contains the server's unique identifier.
// The identifier's scope is best assumed to be bound by the user's account, unless other arrangements have been made with Rackspace.
//
// Any Links provided are used to refer to the server specifically by URL.
// These links are useful for making additional REST calls not explicitly supported by Gorax.
type NewServer struct {
	Name            string          `json:"name",omitempty`
	ImageRef        string          `json:"imageRef,omitempty"`
	FlavorRef       string          `json:"flavorRef,omitempty"`
	OsDcfDiskConfig string          `json:"OS-DCF:diskConfig,omitempty"`
	Metadata        interface{}     `json:"metadata,omitempty"`
	Personality     []FileConfig    `json:"personality,omitempty"`
	Networks        []NetworkConfig `json:"networks,omitempty"`
	AdminPass       string          `json:"adminPass,omitempty"`
	Id              string          `json:"id,omitempty"`
	Links           []Link          `json:"links,omitempty"`
}

// RaxBandwidth provides measurement of server bandwidth consumed over a given audit interval.
type RaxBandwidth struct {
	AuditPeriodEnd    string `json:"audit_period_end"`
	AuditPeriodStart  string `json:"audit_period_start"`
	BandwidthInbound  int64  `json:"bandwidth_inbound"`
	BandwidthOutbound int64  `json:"bandwidth_outbound"`
	Interface         string `json:"interface"`
}

// ResizeRequest structures are used internally to encode to JSON the parameters required to resize a server instance.
// Client applications will not use this structure (no API accepts an instance of this structure).
// See the Region method ResizeServer() for more details on how to resize a server.
type ResizeRequest struct {
	Name       string `json:"name"`
	FlavorRef  string `json:"flavorRef"`
	DiskConfig string `json:"diskConfig"`
}
