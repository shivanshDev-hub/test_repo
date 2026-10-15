# THIS IS A EMPTY FILE 
import requests
import re

class ReconModule:
    def __init__(self, timeout=10):
        self.timeout = timeout
        self.base_url = "https://crt.sh/?q={domain}&output=json"

    def sanitize_input(self, domain):
        """Security: Prevent command injection and malformed inputs."""
        # Only allow alphanumeric, dots, and hyphens
        if not re.match(r"^[a-zA-Z0-9.-]+$", domain):
            raise ValueError("Invalid domain format")
        return domain.lower()

    def fetch_data(self, domain):
        """Integration: Handles the actual network request."""
        sanitized = self.sanitize_input(domain)
        try:
            response = requests.get(self.base_url.format(domain=sanitized), timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            return {"error": str(e)}

    def parse_subdomains(self, json_data):
        """Unit: Logic-only parsing of the response."""
        if not isinstance(json_data, list):
            return []
        
        subdomains = set()
        for entry in json_data:
            name = entry.get("name_value", "")
            # Handle wildcards and multi-line entries
            for sub in name.split('\n'):
                sub = sub.replace("*.", "").strip()
                if sub:
                    subdomains.add(sub)
        return sorted(list(subdomains))
