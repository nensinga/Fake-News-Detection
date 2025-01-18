# Main
import sys  # System functions.
import os  # Operating System functions.
from colorama import Fore  # For text color.

# Importing modules
# SECURITY.
import src.modules.recpull as recpull
# ENUMERATION.
import src.modules.shodan as shodan
import src.modules.numlook as numlook
import src.modules.geolock as geolock
# CASE-GEN.
import src.modules.loki_keygen as loki_keygen
import src.modules.loki_discovery as loki_discovery
import src.modules.loki_encrypt as loki_encrypt
import src.modules.loki_decrypt as loki_decrypt
# OSINT.
import src.modules.cryptotrace as cryptotrace
import src.modules.vt as vt
import src.modules.mactrace as mactrace
import src.modules.ovpn as ovpn
import src.modules.pvpn as pvpn
import src.modules.flightinfo as flightinfo
import src.modules.wigle as wigle
import src.modules.bankindex as bankindex
import src.modules.exif as exif
import src.modules.ytd as ytd
import src.modules.falcon as falcon

from . import apicon
from .utils import print_hero, PROMPT

# Pre-run setup
os.system("clear")

# Hide tracebacks - change to 1 for development mode.
sys.tracebacklimit = 0

# Command-to-function mapping
COMMANDS = {
    "shodan": shodan.run_shodan,
    "numlook": numlook.numlook,
    "geolock": geolock.geolock,
    "cryptotrace": cryptotrace.cryptotrace,
    "vt": vt.vt,
    "ovpn": ovpn.ovpn,
    "pvpn": pvpn.pvpn,
    "mactrace": mactrace.mactrace,
    "flightinfo": flightinfo.flightinfo,
    "wigle": wigle.wigle,
    "bankindex": bankindex.bankindex,
    "ytd": ytd.ytd,
    "lokigen": loki_keygen.loki_keygen,
    "lokidiscovery": loki_discovery.loki_discovery,
    "lokiencrypt": loki_encrypt.loki_encrypt,
    "lokidecrypt": loki_decrypt.loki_decrypt,
    "apicon": apicon.apicon,
    "exif": exif.exif,
    "falcon": falcon.falcon,
    "recpull": recpull.recpull,
}

def main_script():
    try:
        print_hero()
        option = input(f"{PROMPT} ").strip().lower()

        # Execute the corresponding function if the command exists
        if option in COMMANDS:
            COMMANDS[option]()
        else:
            print(f"{Fore.RED}Invalid option: {option}{Fore.WHITE}")

    except KeyboardInterrupt:
        print(f'\n{Fore.YELLOW}You interrupted the program.{Fore.WHITE}')
    finally:
        sys.exit(0)
