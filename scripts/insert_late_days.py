import csv
import requests
import os

TOTAL_LATE_DAYS = 3
API_URL = "http://localhost:8080/v1/late-days"

csv_file_path = os.path.join(os.path.dirname(__file__), '..', 'data', 'MP0.csv')

def process_csv():
    with open(csv_file_path, 'r') as csvfile:
        csvreader = csv.DictReader(csvfile)
        for row in csvreader:
            student_email = row['Student Name']
            student_id = student_email.split('@')[0]
            late_days_used = int(row['Late Days'])
            late_days_left = TOTAL_LATE_DAYS - late_days_used
            if late_days_left < 0:
                late_days_left = -1

            # Prepare the data for the API request
            data = {
                "days": late_days_left,
                "student_email": student_email,
                "student_id": student_id
            }

            # Send POST request to the API
            try:
                response = requests.post(API_URL, json=data)
                response.raise_for_status()
                print(f"Successfully inserted data for {student_email}")
            except requests.exceptions.RequestException as e:
                # print the error message returned by the API
                print(f"Error inserting data for {student_email}: {response.json()['error']}")

if __name__ == "__main__":
    process_csv()