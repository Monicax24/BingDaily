import requests
import threading
import time

"""
TO DO:
- Refactor this script for new API changes
"""

TOKEN = "temp"
URL = "temp"

# Set up the headers with the authentication token
HEADERS = {
    "Authorization": f"Bearer {TOKEN}"
}

# Constants to control the behavior
NUM_THREADS = 100  # Number of threads (each thread simulates one user)
REQUEST_FREQUENCY = 0.5  # Delay in seconds between API calls for each user
RUN_DURATION = 60  # Total duration in seconds to run the simulation

# Shared variables for tracking response times
response_times = []
response_times_lock = threading.Lock()

# Function to send a single GET request
def send_request():
    try:
        start_time = time.time()
        response = requests.get(URL, headers=HEADERS)
        elapsed_time = time.time() - start_time

        # Record the response time
        with response_times_lock:
            response_times.append(elapsed_time)

        print(f"Status Code: {response.status_code}, Response: {response.text[:100]}")  # Print first 100 chars
    except Exception as e:
        print(f"Request failed: {e}")

# Function for each thread to simulate a user
def simulate_user():
    end_time = time.time() + RUN_DURATION
    while time.time() < end_time:
        send_request()
        time.sleep(REQUEST_FREQUENCY)

# Main function to start the simulation
def main():
    threads = []
    for _ in range(NUM_THREADS):
        thread = threading.Thread(target=simulate_user)
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()

    # Calculate and print the average response time
    with response_times_lock:
        if response_times:
            average_response_time = sum(response_times) / len(response_times)
            print(f"Average Response Time: {average_response_time:.2f} seconds")
        else:
            print("No responses recorded.")

if __name__ == "__main__":
    main()
