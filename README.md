# store_monitoring

### Project Details

|                          |                  |
|--------------------------|------------------|
| **Service name**         | store_monitoring |
| **Language / Framework** | Golang/Gin       |
| **Golang Version**       | 1.18             |
| **Database**             | postgreSQL       | 

### CurrentTime
Currently the project has hard-coded currentTime = "2023-01-25 18:13:22.47922 UTC" which max of time present in store_status db. You can edit it in util.go file.
### Database setup in local
1. create a database `store_monitoring` in pgAdmin:

2. We have 5 tables `store_business_hour`, `store_status`, `timezones`, `report_status`, and `report`

3. Commands to create these tables:

    
    CREATE TABLE store_business_hours (
    id SERIAL PRIMARY KEY, store_id bigint NOT NULL,
    day_of_week INTEGER NOT NULL,
    start_time_local TEXT,
    end_time_local TEXT
    );


    CREATE TABLE store_status (                                                       
    id SERIAL PRIMARY KEY, store_id bigint NOT NULL,
    timestamp_utc TIMESTAMP WITH TIME ZONE NOT NULL,
    status text NOT NULL
    );


    CREATE TABLE timezones (
    store_id bigint NOT NULL,
    timezone_str TEXT NOT NULL,
    PRIMARY KEY (store_id)
    );


    CREATE TABLE report_status (
    report_id VARCHAR(32) PRIMARY KEY,
    status VARCHAR(16) NOT NULL
    );


    CREATE TABLE report (
    report_id TEXT,
    store_id BIGINT,
    uptime_last_hour DOUBLE PRECISION,
    uptime_last_day DOUBLE PRECISION,
    uptime_last_week DOUBLE PRECISION,
    downtime_last_hour DOUBLE PRECISION,
    downtime_last_day DOUBLE PRECISION,
    downtime_last_week DOUBLE PRECISION,
    PRIMARY KEY(report_id, store_id)
    );

4. Load CSV files data to `store_business_hour`, `store_status`, and `timezones` tables


    COPY store_business_hours(store_id, day_of_week, start_time_local, end_time_local)
    FROM '/path/business_hours.csv'
    WITH (FORMAT CSV, HEADER);


    COPY store_status(store_id, status, timestamp_utc)
    FROM '/path/store_status.csv'
    WITH (FORMAT CSV, HEADER);


    COPY timezones(store_id, timezone_str)
    FROM '/path/store_timezone.csv'
    WITH (FORMAT CSV, HEADER);

5. Change the `user`, `password`, and `dbname` accordingly in `connectDB()` method in `db.go` 

### API cURLS

    curl --location --request POST 'http://localhost:8080/trigger_report' \
    --header 'Content-Type: application/json'

    Response:
    {
    "ReportId": "akIqYzNL"
    }


    curl --location --request GET 'http://localhost:8080/get_report/:report_id'


### Explanation:

`/trigger_report` API calls `TriggerReportGeneration()`


1. `TriggerReportGeneration()`, this method fetched all the storeIds from `timezone` tables, generates a random string reportId, insert a row in report_status table with values {report_id, "Running"}. Calls a go routine method `triggerReportGenerationForEachStore()` and return the reportId for the response. Instead of calling go routine we could have also used message queue here to trigger report generation for all store Ids.
2. `triggerReportGenerationForEachStore()`, this method calls `GenerateAndStoreReportForStoreId()` for each storeId and once the report is generated for all the storeId, it updated the status of report_id in `report_status` with {report_id, "Completed"}.
3. `GenerateAndStoreReportForStoreId()`, since we have to generate report for last week for each storeId, we are fetching the business_hours of storeId for each each in last week, if data for any day is missing, enrich the data with startTime `00:00:00` and endTime `23:59:59`. Now fetch all the status of this storeId within a week and call `calculateWeeklyObservationAndGererateReport()` to generate the report for this storeId and insert the record in `report` table.
4. `calculateWeeklyObservationAndGererateReport()` I have divided business hours for each into chunks of `60 mins` and initialized each chunk with value `None`. Then iterate over list of status of the store_id and if the time of this status lies in between the business hour of that day, update the chunk in which this lies with the status `inactive` or `active` value. Interpolate the remaining `None` value chunks with `enrichStatusMapWithNearestStatus()` method, generate the weekly report with `createWeeklyObservation()` and return the report.
5. `enrichStatusMapWithNearestStatus()` For each business hour we need to replace `None` status with `active` or `inactive` status. If the day has zero `active` or `inactive` status, I have simply updated the status of each business hour by geerating a random value between 0 to 1, if its >= 0.5, status is `active` else `inactive`. For those days which has some `active` and `inactive` status, then if the status of a chunk is `None` replace it with nearest `active` or `inactive`chunk status. For the starting reference point again I have taken random value generation approach.
6. `createWeeklyObservation()`, now we have `active` or `inactive` status for every chunk of each day, simply used (no. of active chunks / no. of total chunks)*100 to generate uptime percentage of last hour, last day and last week.


`/get_report:report_id` API calls `GetCSVData()`, it fetches status of `report_id` from `report_status` table, if the status is "Running", it returns the status and if the status is "Completed", fetches report for each store_id corresponding to the report_id from `report` table and return the report which later get converted to csv file and response is returned.