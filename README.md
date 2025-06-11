# api-broker
Build the data structures


You need to write a reliable service to return the location (country, city) from a given IP address. Assume that instead of keeping a database of IPs on your side, you use a couple of free third party services (let's call them a provider) to map the given IP to its location (ipinfo.io is one such service). These services being free don't have strict SLAs and can face issues from time to time. We intend to build a broker/proxy in between that routes the request to the service that is most likely going to return a good response at that time. A good response is determined by the following criteria:

1. No error (should be a valid response)
2. Least response time

Also, note that these services have pre-defined thresholds on how many requests can be made per minute and if an additional request is made after the threshold count in that minute then it throws an error.

Expected outcome:

2. Build the data structures for holding quality of service parameters over time per provider (like errors in last 5 mins, avg response time in last 5 mins and requests made in last minute)
3. Write the dynamic routing logic (considering that there will be concurrent requests to the broker/proxy) for sending the request the most appropriate provider
4. Avoid making parallel request across providers to maximize throughput
