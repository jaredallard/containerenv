FROM archlinux/base
ENTRYPOINT ["/init-wrapper"]

# Setup locales
RUN echo -ne "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
&&  locale-gen

# Setup mirrors
RUN pacman -Sy --noconfirm reflector \
&&  reflector -c "US" --sort age --sort rate --threads 100 --save /etc/pacman.d/mirrorlist \
&&  pacman -R --noconfirm reflector

# xorg, pulseaudio, docker, and update
RUN pacman --noconfirm -S archlinux-keyring grep xorg xorg-xinit pulseaudio sudo docker \
&&  pacman -Syu --noconfirm

COPY init /init-wrapper
