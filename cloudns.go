// Package cloudns public structs/functions
package cloudns

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/tidwall/gjson"
)

// Apiaccess ClouDNS API Credentials, see https://www.cloudns.net/wiki/article/42/
type Apiaccess struct {
	Authid       int    `json:"auth-id,omitempty"`
	Subauthid    int    `json:"sub-auth-id,omitempty"`
	Authpassword string `json:"auth-password"`
}

// Zone is the external representation of a zone
// check the ...zone types in api.go for details
type Zone struct {
	Domain string   `json:"domain-name"`
	Ztype  string   `json:"zone-type"`
	Ns     []string `json:"ns,omitempty"`
}

// Listzones returns all zones (max: 100)
func (a Apiaccess) Listzones() ([]Zone, error) {
	zls := zonelist{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Page:         1,
		Hits:         100,
	}
	resp, err := zls.lszone()
	var rz []Zone
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return rz, errors.New(errmsg)
		}
		var intrz []retzone
		err2 := json.Unmarshal(resp.Body(), &intrz)
		for _, zn := range intrz {
			tmpzn := Zone{
				Domain: zn.Domain,
				Ztype:  zn.Ztype,
			}
			rz = append(rz, tmpzn)
		}
		return rz, err2
	}
	return rz, err
}

// List returns all records from a zone
func (z Zone) List(a *Apiaccess) ([]Record, error) {
	var ra []Record
	rls := reclist{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Domain:       z.Domain,
	}
	resp, err := rls.lsrec()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return ra, errors.New(errmsg)
		}
		var ratmp map[string]retrec
		err2 := json.Unmarshal(resp.Body(), &ratmp)
		for _, rec := range ratmp {
			tmpttl, _ := strconv.Atoi(rec.TTL)
			tmppriority, _ := strconv.Atoi(rec.Priority)
			tmpframe, _ := strconv.Atoi(rec.Frame)
			tmpframetitle := rec.FrameTitle
			tmpframekeywords := rec.FrameKeywords
			tmpframedescription := rec.FrameDescription
			tmpmobilemeta := rec.MobileMeta
			tmpsavepath := rec.SavePath
			tmpredirecttype := rec.RedirectType
			tmpweight := rec.Weight
			tmpport := rec.Port
			tmpmail := rec.Mail
			tmptxt := rec.Txt
			tmpalgorithm := rec.Algorithm
			tmpfptype := rec.Fptype
			rectmp := Record{
				Domain:           z.Domain,
				ID:               rec.ID,
				Rtype:            rec.Rtype,
				Host:             rec.Host,
				TTL:              tmpttl,
				Record:           rec.Record,
				Priority:         tmppriority,
				Frame:            tmpframe,
				FrameTitle:       tmpframetitle,
				FrameKeywords:    tmpframekeywords,
				FrameDescription: tmpframedescription,
				MobileMeta:       tmpmobilemeta,
				SavePath:         tmpsavepath,
				RedirectType:     tmpredirecttype,
				Weight:           tmpweight,
				Port:             tmpport,
				Mail:             tmpmail,
				Txt:              tmptxt,
				Algorithm:        tmpalgorithm,
				Fptype:           tmpfptype,
			}
			ra = append(ra, rectmp)
		}
		return ra, err2
	}
	return ra, err
}

// Create a new zone
func (z Zone) Create(a *Apiaccess) (Zone, error) {
	cr := createzone{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Domain:       z.Domain,
		Ztype:        z.Ztype,
		Ns:           z.Ns,
	}
	resp, err := cr.create()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return z, errors.New(errmsg)
		}
	}
	return z, err
}

// Read a zone
func (z Zone) Read(a *Apiaccess) (Zone, error) {
	cr := createzone{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Domain:       z.Domain,
		Ztype:        z.Ztype,
		Ns:           z.Ns,
	}
	resp, err := cr.read()
	var zlint []retzone
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return z, errors.New(errmsg)
		}
		junerr := json.Unmarshal(resp.Body(), &zlint)
		if junerr == nil {
			var rz = Zone{
				Domain: zlint[0].Domain,
				Ztype:  zlint[0].Ztype,
				Ns:     z.Ns,
			}
			return rz, junerr
		}
	}
	return z, err
}

// Update a zone [dummy]
func (z Zone) Update(a *Apiaccess) (Zone, error) {
	err := errors.New("Zone updates are currently not implemented, see https://github.com/sta-travel/cloudns-go/limitations.md")
	return z, err
}

// Destroy a zone
func (z Zone) Destroy(a *Apiaccess) (Zone, error) {
	cr := createzone{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Domain:       z.Domain,
		Ztype:        z.Ztype,
		Ns:           z.Ns,
	}
	resp, err := cr.destroy()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return z, errors.New(errmsg)
		}
	}
	return z, err
}

// Record is the external representation of a record
// check the ...record types in api.go for details
type Record struct {
	ID               string  `json:"id"`
	Domain           string  `json:"domain-name"`
	Host             string  `json:"host"`
	Rtype            string  `json:"record-type"`
	TTL              int     `json:"ttl"`
	Record           string  `json:"record"`
	Priority         int     `json:"priority,omitempty"`
	Weight           int     `json:"weight,omitempty"`
	Port             int     `json:"port,omitempty"`
	Frame            int     `json:"frame,omitempty"`
	FrameTitle       string  `json:"frame-title,omitempty"`
	FrameKeywords    string  `json:"frame-keywords,omitempty"`
	FrameDescription string  `json:"frame-description,omitempty"`
	MobileMeta       int     `json:"mobile-meta,omitempty"`
	SavePath         int     `json:"save-path,omitempty"`
	RedirectType     int     `json:"redirect-type,omitempty"`
	Mail             string  `json:"mail,omitempty"`
	Txt              string  `json:"txt,omitempty"`
	Algorithm        int     `json:"algorithm,omitempty"`
	Fptype           int     `json:"fptype,omitempty"`
	Status           int     `json:"status,omitempty"`
	GeodnsLocation   int     `json:"geodns-location,omitempty"`
	GeodnsCode       string  `json:"geodns-code,omitempty"`
	CaaFlag          int     `json:"caa_flag,omitempty"`
	CaaType          string  `json:"caa_type,omitempty"`
	CaaValue         string  `json:"caa_value,omitempty"`
	TlsaUsage        string  `json:"tlsa_usage,omitempty"`
	TlsaSelector     string  `json:"tlsa_selector,omitempty"`
	TlsaMatchingType string  `json:"tlsa_matching_type,omitempty"`
	KeyTag           int     `json:"key-tag,omitempty"`
	DigestType       int     `json:"digest-type,omitempty"`
	Order            string  `json:"order,omitempty"`
	Pref             string  `json:"pref,omitempty"`
	Flag             int     `json:"flag,omitempty"`
	Params           string  `json:"params,omitempty"`
	Regexp           string  `json:"regexp,omitempty"`
	Replace          int     `json:"replace,omitempty"`
	CertType         int     `json:"cert-type,omitempty"`
	CertKeyTag       int     `json:"cert-key-tag,omitempty"`
	CertAlgorithm    int     `json:"cert-algorithm,omitempty"`
	LatDeg           float32 `json:"lat-deg,omitempty"`
	LatMin           float32 `json:"lat-min,omitempty"`
	LatSec           float32 `json:"lat-sec,omitempty"`
	LatDir           string  `json:"lat-dir,omitempty"`
	LongDeg          float32 `json:"long-deg,omitempty"`
	LongMin          float32 `json:"long-min,omitempty"`
	LongSec          float32 `json:"long-sec,omitempty"`
	LongDir          string  `json:"long-dir,omitempty"`
	Altitude         float32 `json:"altitude,omitempty"`
	Size             float32 `json:"size,omitempty"`
	HPrecision       float32 `json:"h-precision,omitempty"`
	VPrecision       float32 `json:"v-precision,omitempty"`
	CPU              string  `json:"cpu,omitempty"`
	OS               string  `json:"os,omitempty"`
}

// Create a new record
func (r Record) Create(a *Apiaccess) (Record, error) {
	inr := createrec{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Domain:       r.Domain,
		Host:         r.Host,
		Rtype:        r.Rtype,
		TTL:          r.TTL,
		Record:       r.Record,
	}
	if r.Rtype == "MX" || r.Rtype == "SRV" {
		inr.Priority = &r.Priority
	} else if r.Rtype == "WR" {
		inr.Frame = r.Frame
		inr.FrameTitle = r.FrameTitle
		inr.FrameKeywords = r.FrameKeywords
		inr.FrameDescription = r.FrameDescription
		inr.MobileMeta = r.MobileMeta
		inr.SavePath = r.SavePath
		inr.RedirectType = r.RedirectType
	} else if r.Rtype == "SRV" {
		inr.Priority = &r.Priority
		inr.Weight = r.Weight
		inr.Port = r.Port
	} else if r.Rtype == "RP" {
		inr.Mail = r.Mail
		inr.Txt = r.Txt
	}
	resp, err := inr.create()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return r, errors.New(errmsg)
		}
		newid := gjson.GetBytes(resp.Body(), "data.id")
		r.ID = newid.String()
	}
	return r, err
}

// Read a record
func (r Record) Read(a *Apiaccess) (Record, error) {
	lsr := reclist{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Domain:       r.Domain,
		Host:         r.Host,
		Rtype:        r.Rtype,
	}
	resp, err := lsr.lsrec()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return r, errors.New(errmsg)
		}
		var ratmp map[string]retrec
		err2 := json.Unmarshal(resp.Body(), &ratmp)
		for _, rec := range ratmp {
			tmpttl, _ := strconv.Atoi(rec.TTL)
			tmppriority, _ := strconv.Atoi(rec.Priority)
			tmpframe, _ := strconv.Atoi(rec.Frame)
			tmpframetitle := rec.FrameTitle
			tmpframekeywords := rec.FrameKeywords
			tmpframedescription := rec.FrameDescription
			tmpmobilemeta := rec.MobileMeta
			tmpsavepath := rec.SavePath
			tmpredirecttype := rec.RedirectType
			tmpweight := rec.Weight
			tmpport := rec.Port
			tmptxt := rec.Txt
			tmpmail := rec.Mail
			rectmp := Record{
				Domain:           r.Domain,
				ID:               rec.ID,
				Rtype:            rec.Rtype,
				Host:             rec.Host,
				TTL:              tmpttl,
				Record:           rec.Record,
				Priority:         tmppriority,
				Frame:            tmpframe,
				FrameTitle:       tmpframetitle,
				FrameKeywords:    tmpframekeywords,
				FrameDescription: tmpframedescription,
				MobileMeta:       tmpmobilemeta,
				SavePath:         tmpsavepath,
				RedirectType:     tmpredirecttype,
				Weight:           tmpweight,
				Port:             tmpport,
				Txt:              tmptxt,
				Mail:             tmpmail,
			}
			if r.ID != "" && r.ID == rectmp.ID {
				return rectmp, err2
			}
			// if we do not have an ID match, we return the last one ...
			return rectmp, err2
		}
		return r, err2
	}
	return r, err
}

// Update a record
func (r Record) Update(a *Apiaccess) (Record, error) {
	tmpid, _ := strconv.Atoi(r.ID)
	inr := updaterec{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Rid:          tmpid,
		Domain:       r.Domain,
		Host:         r.Host,
		TTL:          r.TTL,
		Record:       r.Record,
	}
	if r.Rtype == "MX" || r.Rtype == "SRV" {
		inr.Priority = &r.Priority
	} else if r.Rtype == "WR" {
		inr.Frame = r.Frame
		inr.FrameTitle = r.FrameTitle
		inr.FrameKeywords = r.FrameKeywords
		inr.FrameDescription = r.FrameDescription
		inr.MobileMeta = r.MobileMeta
		inr.SavePath = r.SavePath
		inr.RedirectType = r.RedirectType
	} else if r.Rtype == "SRV" {
		inr.Priority = &r.Priority
		inr.Weight = r.Weight
		inr.Port = r.Port
	} else if r.Rtype == "RP" {
		inr.Txt = r.Txt
		inr.Mail = r.Mail
	}
	resp, err := inr.update()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return r, errors.New(errmsg)
		}
	}
	return r, err
}

// Destroy a record
func (r Record) Destroy(a *Apiaccess) (Record, error) {
	tmpid, _ := strconv.Atoi(r.ID)
	inr := updaterec{
		Authid:       a.Authid,
		Subauthid:    a.Subauthid,
		Authpassword: a.Authpassword,
		Rid:          tmpid,
		Domain:       r.Domain,
		Host:         r.Host,
		TTL:          r.TTL,
		Record:       r.Record,
	}
	resp, err := inr.destroy()
	if err == nil {
		errmsg, isapierr := checkapierr(resp.Body())
		if isapierr {
			return r, errors.New(errmsg)
		}
	}
	return r, err
}
