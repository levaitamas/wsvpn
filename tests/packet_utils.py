from platform import system
from tests.conftest import GoBin

import scapy.layers.all as scapy_layers
import scapy.plist as scapy_plist
import scapy.packet as scapy_packet
import scapy.sendrecv as scapy_sendrecv

# This is essentially the __eq__ function from Scapy, except it ignores values that are None in either item
def packet_equal(self, other):
    if isinstance(self, scapy_packet.NoPayload):
        return self == other

    if not isinstance(other, self.__class__):
        return False

    for f in self.fields_desc:
        if f not in other.fields_desc:
            return False
        
        self_val = self.getfieldval(f.name)
        other_val = other.getfieldval(f.name)

        if self_val is not None and other_val is not None and self_val != other_val:
            return False

    return packet_equal(self.payload, other.payload)

class PacketTest:
    def __init__(self, svbin: GoBin, clbin: GoBin, ethernet: bool) -> None:
        self.svbin = svbin
        self.clbin = clbin
        self.ethernet = ethernet
        self.pkts = []
        self.need_dummy_layer = (not ethernet) and (system() == "Darwin")

    def pkt_add(self, pkt):
        if self.ethernet:
            pkt = scapy_layers.Ether()/pkt
        self.pkts.append(pkt)

    def simple_pkt(self, pktlen: int):
        payload = scapy_layers.ICMP(type=0, code=0, id=0x0, seq=0x0, chksum=0xf7ff)
        if pktlen > 0:
            payload = payload / scapy_packet.Raw(bytes(b"A"*pktlen))
        
        pkt = scapy_layers.IP(version=4,ihl=5,src="192.168.3.2",dst="192.168.3.1",len=28+pktlen,id=len(self.pkts)+1,ttl=1) / payload

        if self.need_dummy_layer:
            pkt = scapy_layers.Loopback(type=0x2) / pkt

        self.pkt_add((pkt, pkt))

    def add_defaults(self):
        self.simple_pkt(0)
        self.simple_pkt(10)
        self.simple_pkt(1000)
        self.simple_pkt(1300)

    def run(self):
        self.svbin.assert_ready_ok()
        self.clbin.assert_ready_ok()

        for pkt, raw_pkt in self.pkts:
            send_iface = None
            recv_iface = None

            def sendpkt():
                scapy_sendrecv.sendp(raw_pkt, iface=send_iface, count=1, return_packets=True)

            def dosniff() -> scapy_plist.PacketList:
                res: scapy_plist.PacketList = scapy_sendrecv.sniff(iface=recv_iface, started_callback=sendpkt, filter=None, count=1, store=1, timeout=2)
                assert len(res.res) > 0

                actual_pkt = res.res[0]

                assert packet_equal(pkt, actual_pkt)

            send_iface = "utun100"
            recv_iface = "utun20"
            dosniff()

            send_iface = "utun20"
            recv_iface = "utun100"
            dosniff()