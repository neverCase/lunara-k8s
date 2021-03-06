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
    && docker-php-ext-install pcntl \
    && pecl install redis-4.3.0 swoole-1.9.23 \
    && docker-php-ext-enable redis.so swoole.so \
    && apt-get purge -y --auto-remove $buildDeps \
    && rm -r /var/lib/apt/lists/* \
#    install lua-5.3.4 include liblua.so
    && wget http://www.lua.org/ftp/lua-5.3.4.tar.gz \
    && tar -zvxf lua-5.3.4.tar.gz \
    && cd lua-5.3.4 \
    && sed -i 's/TO_LIB= liblua.a/TO_LIB= liblua.a liblua.so/g' Makefile \
    && sed -i 's/$(MYCFLAGS)/$(MYCFLAGS) -fPIC/g' src/Makefile \
    && sed -i 's/LUA_O=\tlua.o/LUA_O=\tlua.o\n\nLUA_SO=\tliblua.so/g' src/Makefile \
    && sed -i 's/ALL_T= $(LUA_A) $(LUA_T) $(LUAC_T)/ALL_T= $(LUA_A) $(LUA_T) $(LUAC_T) $(LUA_SO)/g' src/Makefile \
    && sed -i 's/clean:/$(LUA_SO): $(CORE_O) $(LIB_O)\n\t$(CC) -o $@ -shared $? -ldl -lm\n\nclean:/g' src/Makefile \
    && make linux test \
    && make install \
    && cp src/lua /usr/bin/lua \
    && cp src/luac /usr/bin/luac \
    && cp src/liblua.so /usr/lib/x86_64-linux-gnu/liblua5.3.so \
    && ln -s /usr/local/lib/liblua.a /usr/lib/liblua.a \
    && mkdir /usr/include/lua \
    && ln -s /usr/local/include/*lua* /usr/include/lua/ \
    && cd .. \
    && rm -rf lua-5.3.4.tar.gz \
    && rm -rf lua-5.3.4 \
#    install php lua.so
    && wget http://pecl.php.net/get/lua-1.0.0.tgz \
    && tar -zvxf lua-1.0.0.tgz \
    && cd lua-1.0.0 \
    && phpize \
    && ./configure --with-php-config=/usr/local/bin/php-config \
    && make \
    && make install \
    && docker-php-ext-enable lua.so \
    && cd .. \
    && rm -rf lua-1.0.0.tgz \
    && rm -rf lua-1.0.0 \
#    install libevent-dev
    && wget https://github.com/libevent/libevent/releases/download/release-2.1.11-stable/libevent-2.1.11-stable.tar.gz \
    && tar -zvxf libevent-2.1.11-stable.tar.gz \
    && cd libevent-2.1.11-stable \
    && ./configure --prefix=/usr/local/libevent-2.1.11/ \
    && make \
    && make install \
    && cd .. \
    && rm -rf libevent-2.1.11-stable.tar.gz \
    && rm -rf libevent-2.1.11-stable \
#    install php libevent.so
    && wget https://pecl.php.net/get/libevent-0.1.0.tgz \
    && tar -zxvf libevent-0.1.0.tgz \
    && cd libevent-0.1.0 \
    && phpize \
    && ./configure --with-libevent=/usr/local/libevent-2.1.11/ \
    && make \
    && make install \
    && docker-php-ext-enable libevent.so \
    && cd .. \
    && rm -rf libevent-0.1.0.tgz \
    && rm -rf libevent-0.1.0

WORKDIR /data/
COPY daemon-hook /data/

CMD ["/data/daemon-hook"]
