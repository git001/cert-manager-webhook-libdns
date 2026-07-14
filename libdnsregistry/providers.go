// This file is auto generated, DO NOT EDIT.
package libdnsregistry

import (
	libdnsacmedns "github.com/libdns/acmedns"
	libdnsacmeproxy "github.com/libdns/acmeproxy"
	libdnsalidns "github.com/libdns/alidns"
	libdnsall_inkl "github.com/libdns/all-inkl"
	libdnsarvancloud "github.com/libdns/arvancloud"
	libdnsautodns "github.com/libdns/autodns"
	libdnsazure "github.com/libdns/azure"
	libdnsbluecat "github.com/libdns/bluecat"
	libdnsbunny "github.com/libdns/bunny"
	libdnscloudflare "github.com/libdns/cloudflare"
	libdnscloudns "github.com/libdns/cloudns"
	libdnsconoha "github.com/libdns/conoha"
	libdnsdesec "github.com/libdns/desec"
	libdnsdigitalocean "github.com/libdns/digitalocean"
	libdnsdirectadmin "github.com/libdns/directadmin"
	libdnsdnsexit "github.com/libdns/dnsexit"
	libdnsdnsimple "github.com/libdns/dnsimple"
	libdnsdnsupdate "github.com/libdns/dnsupdate"
	libdnsdomainnameshop "github.com/libdns/domainnameshop"
	libdnsduckdns "github.com/libdns/duckdns"
	libdnsdynu "github.com/libdns/dynu"
	libdnsdynv6 "github.com/libdns/dynv6"
	libdnseasydns "github.com/libdns/easydns"
	libdnsedgeone "github.com/libdns/edgeone"
	libdnsgandi "github.com/libdns/gandi"
	libdnsgcore "github.com/libdns/gcore"
	libdnsglesys "github.com/libdns/glesys"
	libdnsgodaddy "github.com/libdns/godaddy"
	libdnsgoogleclouddns "github.com/libdns/googleclouddns"
	libdnshe "github.com/libdns/he"
	libdnshetzner "github.com/libdns/hetzner/v2"
	libdnshuaweicloud "github.com/libdns/huaweicloud"
	libdnsinfomaniak "github.com/libdns/infomaniak"
	libdnsinwx "github.com/libdns/inwx"
	libdnsionos "github.com/libdns/ionos"
	libdnslinode "github.com/libdns/linode"
	libdnsloopia "github.com/libdns/loopia"
	libdnsluadns "github.com/libdns/luadns"
	libdnsmailinabox "github.com/libdns/mailinabox"
	libdnsmetaname "github.com/libdns/metaname"
	libdnsmijnhost "github.com/libdns/mijnhost"
	libdnsmythicbeasts "github.com/libdns/mythicbeasts"
	libdnsnamecheap "github.com/libdns/namecheap"
	libdnsnamesilo "github.com/libdns/namesilo"
	libdnsnetcup "github.com/libdns/netcup"
	libdnsnetlify "github.com/libdns/netlify"
	libdnsnetnod "github.com/libdns/netnod"
	libdnsnfsn "github.com/libdns/nfsn"
	libdnsnjalla "github.com/libdns/njalla"
	libdnsoraclecloud "github.com/libdns/oraclecloud"
	libdnsovh "github.com/libdns/ovh"
	libdnsporkbun "github.com/libdns/porkbun"
	libdnspowerdns "github.com/libdns/powerdns"
	libdnspph "github.com/libdns/pph"
	libdnsregery "github.com/libdns/regery"
	libdnsregfish "github.com/libdns/regfish"
	libdnsrfc2136 "github.com/libdns/rfc2136"
	libdnsroute53 "github.com/libdns/route53"
	libdnsscaleway "github.com/libdns/scaleway"
	libdnssimplydotcom "github.com/libdns/simplydotcom"
	libdnsspaceship "github.com/libdns/spaceship"
	libdnstecnocratica "github.com/libdns/tecnocratica"
	libdnstemplate "github.com/libdns/template"
	libdnstencentcloud "github.com/libdns/tencentcloud"
	libdnsthelittlehost "github.com/libdns/thelittlehost"
	libdnstimeweb "github.com/libdns/timeweb"
	libdnstransip "github.com/libdns/transip"
	libdnsunifi "github.com/libdns/unifi"
	libdnsvolcengine "github.com/libdns/volcengine"
	libdnsvultr "github.com/libdns/vultr/v2"
	libdnswedos "github.com/libdns/wedos"
	libdnswestcn "github.com/libdns/westcn"
)

var registry = RegistryStore{
	"acmedns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsacmedns.Provider](conf)
		},
	},
	"acmeproxy": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsacmeproxy.Provider](conf)
		},
	},
	"alidns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsalidns.Provider](conf)
		},
	},
	"all-inkl": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsall_inkl.Provider](conf)
		},
	},
	"arvancloud": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsarvancloud.Provider](conf)
		},
	},
	"autodns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsautodns.Provider](conf)
		},
	},
	"azure": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsazure.Provider](conf)
		},
	},
	"bluecat": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsbluecat.Provider](conf)
		},
	},
	"bunny": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsbunny.Provider](conf)
		},
	},
	"cloudflare": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnscloudflare.Provider](conf)
		},
	},
	"cloudns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnscloudns.Provider](conf)
		},
	},
	"conoha": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsconoha.Provider](conf)
		},
	},
	"desec": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdesec.Provider](conf)
		},
	},
	"digitalocean": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdigitalocean.Provider](conf)
		},
	},
	"directadmin": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdirectadmin.Provider](conf)
		},
	},
	"dnsexit": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsexit.Provider](conf)
		},
	},
	"dnsimple": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsimple.Provider](conf)
		},
	},
	"dnsupdate": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsupdate.Provider](conf)
		},
	},
	"domainnameshop": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdomainnameshop.Provider](conf)
		},
	},
	"duckdns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsduckdns.Provider](conf)
		},
	},
	"dynu": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynu.Provider](conf)
		},
	},
	"dynv6": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynv6.Provider](conf)
		},
	},
	"easydns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnseasydns.Provider](conf)
		},
	},
	"edgeone": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsedgeone.Provider](conf)
		},
	},
	"gandi": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsgandi.Provider](conf)
		},
	},
	"gcore": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsgcore.Provider](conf)
		},
	},
	"glesys": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsglesys.Provider](conf)
		},
	},
	"godaddy": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsgodaddy.Provider](conf)
		},
	},
	"googleclouddns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsgoogleclouddns.Provider](conf)
		},
	},
	"he": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshe.Provider](conf)
		},
	},
	"hetzner": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshetzner.Provider](conf)
		},
	},
	"huaweicloud": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshuaweicloud.Provider](conf)
		},
	},
	"infomaniak": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinfomaniak.Provider](conf)
		},
	},
	"inwx": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinwx.Provider](conf)
		},
	},
	"ionos": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsionos.Provider](conf)
		},
	},
	"linode": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnslinode.Provider](conf)
		},
	},
	"loopia": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsloopia.Provider](conf)
		},
	},
	"luadns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsluadns.Provider](conf)
		},
	},
	"mailinabox": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmailinabox.Provider](conf)
		},
	},
	"metaname": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmetaname.Provider](conf)
		},
	},
	"mijnhost": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmijnhost.Provider](conf)
		},
	},
	"mythicbeasts": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmythicbeasts.Provider](conf)
		},
	},
	"namecheap": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamecheap.Provider](conf)
		},
	},
	"namesilo": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamesilo.Provider](conf)
		},
	},
	"netcup": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnetcup.Provider](conf)
		},
	},
	"netlify": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnetlify.Provider](conf)
		},
	},
	"netnod": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnetnod.Provider](conf)
		},
	},
	"nfsn": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnfsn.Provider](conf)
		},
	},
	"njalla": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnjalla.Provider](conf)
		},
	},
	"oraclecloud": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsoraclecloud.Provider](conf)
		},
	},
	"ovh": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsovh.Provider](conf)
		},
	},
	"porkbun": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsporkbun.Provider](conf)
		},
	},
	"powerdns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnspowerdns.Provider](conf)
		},
	},
	"pph": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnspph.Provider](conf)
		},
	},
	"regery": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsregery.Provider](conf)
		},
	},
	"regfish": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsregfish.Provider](conf)
		},
	},
	"rfc2136": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsrfc2136.Provider](conf)
		},
	},
	"route53": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsroute53.Provider](conf)
		},
	},
	"scaleway": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsscaleway.Provider](conf)
		},
	},
	"simplydotcom": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnssimplydotcom.Provider](conf)
		},
	},
	"spaceship": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsspaceship.Provider](conf)
		},
	},
	"tecnocratica": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstecnocratica.Provider](conf)
		},
	},
	"template": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstemplate.Provider](conf)
		},
	},
	"tencentcloud": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstencentcloud.Provider](conf)
		},
	},
	"thelittlehost": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsthelittlehost.Provider](conf)
		},
	},
	"timeweb": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstimeweb.Provider](conf)
		},
	},
	"transip": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstransip.Provider](conf)
		},
	},
	"unifi": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsunifi.Provider](conf)
		},
	},
	"volcengine": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsvolcengine.Provider](conf)
		},
	},
	"vultr": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsvultr.Provider](conf)
		},
	},
	"wedos": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnswedos.Provider](conf)
		},
	},
	"westcn": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnswestcn.Provider](conf)
		},
	},
}
