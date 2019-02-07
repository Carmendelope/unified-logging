#
# Copyright (C) 2019 Nalej Group - All Rights Reserved
#

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

# Now just add the binary\
RUN mkdir /nalej

# Now just add the binary
COPY unified-logging-coord /nalej/

ENTRYPOINT ["./nalej/unified-logging-coord"]
