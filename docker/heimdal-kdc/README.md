# Heimdal KDC

Extremely basic Heimdal KDC image.

## Configuration

### Environment Variables

- `REALM`: The realm for the kdc (all-caps)
- `KDC_ADDRESS`: The FQDN for the kdc's realm

### Data

- Data goes in `/data` (should be a volume).
- The heimdal master key (see below) should be mounted at `/secrets/heimdal.mkey`.

### Creating a Master Key

- `docker run --entrypoint=genmaster <thisimage>` (outputs as base64)
