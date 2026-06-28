FROM docker.io/postgis/postgis:18-3.6

LABEL description="PostGIS with H3 extension for geospatial applications"

# "-y" automatically confirms install (Docker build is non interactive session) 
# "--no-install-recommends" not install opional dependencies (smaller image size)
# then clean chache after installation
RUN apt-get update && apt-get install -y --no-install-recommends postgresql-18-h3 && \
	rm -rf /var/lib/apt/lists/*

# 3. (Optional) Copy initialization script to create the extension automatically
# COPY ./docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/
# RUN chmod 755 /docker-entrypoint-initdb.d/*.sh 2>/dev/null || true
