package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

func extractTcmsgAttributes(action int, data []byte, info *Attribute) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var options []byte
	var xStats []byte
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaKind:
			info.Kind = ad.String()
		case tcaOptions:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is known at this moment,
			// so we save it for later
			options = ad.Bytes()
		case tcaChain:
			info.Chain = uint32Ptr(ad.Uint32())
		case tcaXstats:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is know at this moment,
			// so we save it for later
			xStats = ad.Bytes()
		case tcaStats:
			tcstats := &Stats{}
			err := unmarshalStruct(ad.Bytes(), tcstats)
			concatError(multiError, err)
			info.Stats = tcstats
		case tcaStats2:
			tcstats2 := &Stats2{}
			err := unmarshalStruct(ad.Bytes(), tcstats2)
			concatError(multiError, err)
			info.Stats2 = tcstats2
		case tcaHwOffload:
			info.HwOffload = uint8Ptr(ad.Uint8())
		case tcaEgressBlock:
			info.EgressBlock = uint32Ptr(ad.Uint32())
		case tcaIngressBlock:
			info.IngressBlock = uint32Ptr(ad.Uint32())
		case tcaStab:
			stab := &Stab{}
			err := unmarshalStab(ad.Bytes(), stab)
			concatError(multiError, err)
			info.Stab = stab
		default:
			return fmt.Errorf("extractTcmsgAttributes()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	concatError(multiError, ad.Err())
	if multiError != nil {
		return err
	}

	if len(options) > 0 {
		if (action&actionMask == actionQdisc) && hasQOpt(info.Kind) {
			err = extractQOpt(options, info, info.Kind)
		} else {
			err = extractTCAOptions(options, info, info.Kind)
		}
		if err != nil {
			return err
		}
	}

	if len(xStats) > 0 {
		tcxstats := &XStats{}
		if err := extractXStats(xStats, tcxstats, info.Kind); err != nil {
			return err
		}
		info.XStats = tcxstats
	}
	return nil
}

func hasQOpt(kind string) bool {
	classful := map[string]bool{
		"hfsc": true,
		"qfq":  true,
		"htb":  true,
	}
	if _, ok := classful[kind]; ok {
		return true
	}
	return false
}

func extractQOpt(data []byte, tc *Attribute, kind string) error {
	var multiError error
	switch kind {
	case "hfsc":
		info := &HfscQOpt{}
		err := unmarshalHfscQOpt(data, info)
		concatError(multiError, err)
		tc.HfscQOpt = info
	case "qfq":
		info := &Qfq{}
		err := unmarshalQfq(data, info)
		concatError(multiError, err)
		tc.Qfq = info
	case "htb":
		info := &Htb{}
		err := unmarshalHtb(data, info)
		concatError(multiError, err)
		tc.Htb = info
	default:
		return fmt.Errorf("no QOpts for %s", kind)
	}
	return multiError
}

func extractTCAOptions(data []byte, tc *Attribute, kind string) error {
	var multiError error
	switch kind {
	case "choke":
		info := &Choke{}
		err := unmarshalChoke(data, info)
		concatError(multiError, err)
		tc.Choke = info
	case "fq_codel":
		info := &FqCodel{}
		err := unmarshalFqCodel(data, info)
		concatError(multiError, err)
		tc.FqCodel = info
	case "codel":
		info := &Codel{}
		err := unmarshalCodel(data, info)
		concatError(multiError, err)
		tc.Codel = info
	case "fq":
		info := &Fq{}
		err := unmarshalFq(data, info)
		concatError(multiError, err)
		tc.Fq = info
	case "pie":
		info := &Pie{}
		err := unmarshalPie(data, info)
		concatError(multiError, err)
		tc.Pie = info
	case "hhf":
		info := &Hhf{}
		err := unmarshalHhf(data, info)
		concatError(multiError, err)
		tc.Hhf = info
	case "htb":
		info := &Htb{}
		err := unmarshalHtb(data, info)
		concatError(multiError, err)
		tc.Htb = info
	case "hfsc":
		info := &Hfsc{}
		err := unmarshalHfsc(data, info)
		concatError(multiError, err)
		tc.Hfsc = info
	case "dsmark":
		info := &Dsmark{}
		err := unmarshalDsmark(data, info)
		concatError(multiError, err)
		tc.Dsmark = info
	case "drr":
		info := &Drr{}
		err := unmarshalDrr(data, info)
		concatError(multiError, err)
		tc.Drr = info
	case "cbq":
		info := &Cbq{}
		err := unmarshalCbq(data, info)
		concatError(multiError, err)
		tc.Cbq = info
	case "atm":
		info := &Atm{}
		err := unmarshalAtm(data, info)
		concatError(multiError, err)
		tc.Atm = info
	case "pfifo_fast":
		fallthrough
	case "prio":
		info := &Prio{}
		err := unmarshalPrio(data, info)
		concatError(multiError, err)
		tc.Prio = info
	case "tbf":
		info := &Tbf{}
		err := unmarshalTbf(data, info)
		concatError(multiError, err)
		tc.Tbf = info
	case "sfb":
		info := &Sfb{}
		err := unmarshalSfb(data, info)
		concatError(multiError, err)
		tc.Sfb = info
	case "sfq":
		info := &Sfq{}
		err := unmarshalSfq(data, info)
		concatError(multiError, err)
		tc.Sfq = info
	case "red":
		info := &Red{}
		err := unmarshalRed(data, info)
		concatError(multiError, err)
		tc.Red = info
	case "pfifo":
		limit := &FifoOpt{}
		err := unmarshalStruct(data, limit)
		concatError(multiError, err)
		tc.Pfifo = limit
	case "mqprio":
		info := &MqPrio{}
		err := unmarshalMqPrio(data, info)
		concatError(multiError, err)
		tc.MqPrio = info
	case "bfifo":
		limit := &FifoOpt{}
		err := unmarshalStruct(data, limit)
		concatError(multiError, err)
		tc.Bfifo = limit
	case "clsact":
		return extractClsact(data)
	case "ingress":
		return extractIngress(data)
	case "qfq":
		info := &Qfq{}
		err := unmarshalQfq(data, info)
		concatError(multiError, err)
		tc.Qfq = info
	case "basic":
		info := &Basic{}
		err := unmarshalBasic(data, info)
		concatError(multiError, err)
		tc.Basic = info
	case "bpf":
		info := &Bpf{}
		err := unmarshalBpf(data, info)
		concatError(multiError, err)
		tc.BPF = info
	case "cgroup":
		info := &Cgroup{}
		err := unmarshalCgroup(data, info)
		concatError(multiError, err)
		tc.Cgroup = info
	case "u32":
		info := &U32{}
		err := unmarshalU32(data, info)
		concatError(multiError, err)
		tc.U32 = info
	case "flower":
		info := &Flower{}
		err := unmarshalFlower(data, info)
		concatError(multiError, err)
		tc.Flower = info
	case "rsvp":
		info := &Rsvp{}
		err := unmarshalRsvp(data, info)
		concatError(multiError, err)
		tc.Rsvp = info
	case "route4":
		info := &Route4{}
		err := unmarshalRoute4(data, info)
		concatError(multiError, err)
		tc.Route4 = info
	case "fw":
		info := &Fw{}
		err := unmarshalFw(data, info)
		concatError(multiError, err)
		tc.Fw = info
	case "flow":
		info := &Flow{}
		err := unmarshalFlow(data, info)
		concatError(multiError, err)
		tc.Flow = info
	case "matchall":
		info := &Matchall{}
		err := unmarshalMatchall(data, info)
		concatError(multiError, err)
		tc.Matchall = info
	case "netem":
		info := &Netem{}
		err := unmarshalNetem(data, info)
		concatError(multiError, err)
		tc.Netem = info
	case "cake":
		info := &Cake{}
		err := unmarshalCake(data, info)
		concatError(multiError, err)
		tc.Cake = info
	default:
		return fmt.Errorf("extractTCAOptions(): unsupported kind %s: %w", kind, ErrUnknownKind)
	}

	return multiError
}

func extractXStats(data []byte, tc *XStats, kind string) error {
	var multiError error
	switch kind {
	case "sfb":
		info := &SfbXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Sfb = info
	case "red":
		info := &RedXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Red = info
	case "choke":
		info := &ChokeXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Choke = info
	case "htb":
		info := &HtbXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Htb = info
	case "cbq":
		info := &CbqXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Cbq = info
	case "codel":
		info := &CodelXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Codel = info
	case "hhf":
		info := &HhfXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Hhf = info
	case "pie":
		info := &PieXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Pie = info
	case "fq_codel":
		info := &FqCodelXStats{}
		err := unmarshalFqCodelXStats(data, info)
		concatError(multiError, err)
		tc.FqCodel = info
	case "hfsc":
		info := &HfscXStats{}
		err := unmarshalStruct(data, info)
		concatError(multiError, err)
		tc.Hfsc = info
	default:
		return fmt.Errorf("extractXStats(): unsupported kind: %s", kind)
	}
	return multiError
}

func extractClsact(data []byte) error {
	// Clsact is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("clsact is parameterless: %w", ErrInvalidArg)
	}
	return nil
}

func extractIngress(data []byte) error {
	// Ingress is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("extractIngress()\t%v", data)
	}
	return nil
}

const (
	tcaUnspec = iota
	tcaKind
	tcaOptions
	tcaStats
	tcaXstats
	tcaRate
	tcaFcnt
	tcaStats2
	tcaStab
	tcaPad
	tcaDumpInvisible
	tcaChain
	tcaHwOffload
	tcaIngressBlock
	tcaEgressBlock
)
