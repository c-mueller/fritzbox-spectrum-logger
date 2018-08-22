# Fritz!Box Spectrum Logger

Fritz!Box Spectrum Logger (short FSL) is a Utility to periodically download the DSL Spectrum of a Fritz!Box.

## Setup

### Using a Configuration File

#### Example Configuration

```yaml
credentials:
  endpoint: 192.168.178.1               # The IP Address of the Fritz!Box
  username: ""                          # The username used to login (leave this blank if you don't use a username to login)
  password: mypassw0rd                  # The password used to login (leave blank if your Fritz!Box is not protected)
database_path: spectra.db               # The path to store the BoltDB Database at (Absolute or relative)
update_interval: 60                     # Th interval in seconds at wich the spectrum should be Downloaded
autolaunch: false                       # Set to true to automatically start logging once the application starts
bind_address: :8080                     # The address at wich the application should listen on
session_refresh_interval: 3600          # The time, in seconds, between session renewals
session_refresh_attempts: 5             # The count of login attempts that can fail before stopping the logging routine 
max_download_fails: 480                 # The count of failed downloads needed to stop logging 
```

**Notes**:

- `credentials/endpoint`: DNS names might work, have not been tested though
- `credentials/username` and `credentials/password`: Leave both empty if yout FB is not protected, only enter the password if you use Password only authentication.
- `bind_address`: Golang notation for a binding address `:8080` binds to port 8080, accepting from any interface. `127.0.0.1:8080` also binds to port 8080 but only accepts from Localhost.
- `session_refresh_attempts`: The intention here is to allow a restart of the Fritz!Box, without causing the logger to stop. See Logging Routine for more info.


### Using Environment Variables

## Tested Devices

- Fritz!Box 7390 (Firmware Version 6.83)

Other devices like the Fritz!Box 7490 should also work, but they have not been tested!
If you can confirm the functionality of this application on another Fritz!Box, please modify this document and submit a Pull Request.

## Building FSL

## Logging Routine

## License

FSL is licensed under the GPL v3.0 License.

