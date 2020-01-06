FROM centos:latest
FROM php:5.6

RUN buildDeps=" \
        default-libmysqlclient-dev \
        libbz2-dev \
        libmemcached-dev \
        libsasl2-dev \
    " \
    runtimeDeps=" \
            curl \
            wget \
            git \
            libfreetype6-dev \
            libicu-dev \
            libjpeg-dev \
            libldap2-dev \
            libmcrypt-dev \
            libmemcachedutil2 \
            libpng-dev \
            libpq-dev \
            libxml2-dev \
            libreadline-dev \
            gcc \
        " \
    && apt-get update && apt-get install -y $buildDeps $runtimeDeps \
    && docker-php-ext-install bcmath bz2 calendar iconv intl mbstring mcrypt mysql mysqli opcache pdo_mysql pdo_pgsql pgsql soap zip \
    && docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ \
    && docker-php-ext-install gd \
    && docker-php-ext-configure ldap --with-libdir=lib/x86_64-linux-gnu/ \
    && docker-php-ext-install ldap \
    && docker-php-ext-install exif \
    && pecl install memcached-2.2.0 redis-4.3.0 \
    && docker-php-ext-enable memcached.so redis.so \
    && apt-get purge -y --auto-remove $buildDeps \
    && rm -r /var/lib/apt/lists/*
#    install lua-5.3.4
    && wget http://www.lua.org/ftp/lua-5.3.4.tar.gz \
    && tar -zvxf lua-5.3.4.tar.gz \
    && cd lua-5.3.4 \
    && sed -i 's/TO_LIB= liblua.a/TO_LIB= liblua.a liblua.so/g' Makefile
#    && sed -i 's/TO_LIB\= liblua\.a/TO_LIB\= liblua\.a liblua\.so/g' Makefile \
    && sed -i 's/$(MYCFLAGS)/$(MYCFLAGS) -fPIC/g' src/Makefile \
    && sed -i 's/ALL_T= $(LUA_A) $(LUA_T) $(LUAC_T)/ALL_T= $(LUA_A) $(LUA_T) $(LUAC_T) $(LUA_SO)/g' src/Makefile \
    && sed -i 's/LUA_O=  lua.o\\n/LUA_O=  lua.o\\n/g' src/Makefile \
    && make linux test \
    && make install \
    && cp src/lua /usr/bin/lua \
    && cp src/luac /usr/bin/luac \
    && cp src/liblua.so /usr/lib/x86_64-linux-gnu/liblua5.3.so

WORKDIR /data/
COPY daemon-hook /data/

CMD ["/data/daemon-hook"]
