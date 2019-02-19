#
# Copyright (C) 2019 Nalej Group - All Rights Reserved
#

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

RUN mkdir /nalej

# Put certificate in place
RUN wget https://letsencrypt.org/certs/letsencryptauthorityx3.pem.txt -O /nalej/cert.pem

# Now just add the binary
COPY unified-logging-coord /nalej/

ENTRYPOINT ["./nalej/unified-logging-coord"]
