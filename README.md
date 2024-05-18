# License API mock datafeed

## Build and Run 
``` sh 
make d.build && make up
```

## API 

### GET /musicians/payout
**Params:** 
* musician_id (uuid)   
* musician_name (string)
* year (int)
* month (int)

**Example:**
``` bash
curl --location 'http://localhost:3112/musicians/payout?musician_name=bladee&year=2024&month=2'
```
**Response:**
Content-Type text/plain
```
49
```
