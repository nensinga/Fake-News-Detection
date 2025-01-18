# Main Script
import sys  # System functions
import os  # Operating System functions
from colorama import Fore  # For text color

# Importing modules
from src.modules import (
    recpull,
    shodan,
    numlook,
    geolock,
    loki_keygen,
    loki_discovery,
    loki_encrypt,
    loki_decrypt,
    cryptotrace,
    vt,
    mactrace,
    ovpn,
    pvpn,
    flightinfo,
    wigle,
    bankindex,
    exif,
    ytd,
    falcon,
)
from . import apicon
from .utils import print_hero, PROMPT

# Pre-run setup
os.system("clear")

# Set traceback limit (0 = hide tracebacks, 1 = show in dev mode)
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

def execute_command(option):
    """
    Execute the corresponding function for the given option.
    """
    command_function = COMMANDS.get(option)
    if command_function:
        command_function()
    else:
        print(f"{Fore.RED}Invalid option: {option}{Fore.WHITE}")

def main_script():
    """
    Main entry point for the script.
    """
    try:
        print_hero()
        option = input(f"{PROMPT} ").strip().lower()
        execute_command(option)
    except KeyboardInterrupt:
        print(f'\n{Fore.YELLOW}You interrupted the program.{Fore.WHITE}')
    finally:
        sys.exit(0)

# Run the script if executed directly
if __name__ == "__main__":
    main_script()
