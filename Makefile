# DEBOS_OPTS can be overridden with:
#     make DEBOS_OPTS=... all
# USE_CONTAINER can be set to yes/no/auto (default: auto)
#     make USE_CONTAINER=yes all    # Force container use
#     make USE_CONTAINER=no all     # Force native debos

# To build large images, the debos resource defaults are not sufficient. These
# provide defaults that work for us as universally as we can manage.
FAKEMACHINE_BACKEND = $(shell [ -c /dev/kvm ] && echo kvm || echo qemu)
DEBOS_OPTS := --fakemachine-backend $(FAKEMACHINE_BACKEND) --memory 1GiB --scratchsize 6GiB
DEBOS := debos $(DEBOS_OPTS)

# Use http_proxy from the environment, or apt's http_proxy if set, to speed up
# builds.
http_proxy ?= $(shell apt-config dump --format '%v%n' Acquire::http::Proxy)
export http_proxy

.PHONY: all
all: disk-ufs.img disk-sdcard.img

rootfs.tar dtbs.tar.gz: debos-recipes/qualcomm-linux-debian-rootfs.yaml
	$(DEBOS_CMD) $<

DISK_UFS_IMAGES := disk-ufs.img \
	disk-ufs.img1 \
	disk-ufs.img2

$(DISK_UFS_IMAGES): debos-recipes/qualcomm-linux-debian-image.yaml rootfs.tar
	$(DEBOS_CMD) $<

DISK_SDCARD_IMAGES := disk-sdcard.img \
	disk-sdcard.img1 \
	disk-sdcard.img2

$(DISK_SDCARD_IMAGES): debos-recipes/qualcomm-linux-debian-image.yaml rootfs.tar
	$(DEBOS_CMD) -t imagetype:sdcard $<

.PHONY: test
test: disk-ufs.img
	# rootfs/ is a build artifact, so should not be scanned for tests
	py.test-3 --ignore=rootfs

.PHONY: clean
clean:
	rm -f $(DISK_UFS_IMAGES)
	rm -f $(DISK_SDCARD_IMAGES)
	rm -f rootfs.tar
	rm -f dtbs.tar.gz
