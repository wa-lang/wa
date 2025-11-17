# Firmware Image Format

https://docs.espressif.com/projects/esptool/en/latest/esp32c3/advanced-topics/firmware-image-format.html

This is technical documentation for the firmware image format used by the ROM bootloader. These are the images created by `esptool elf2image`.

![](./images/firmware_image_format.png)

The firmware file consists of a header, an extended header, a variable number of data segments and a footer. Multi-byte fields are little-endian.

## File Header

![](./images/firmware_image_header_format.png)

The image header is 8 bytes long:

| Byte | Description |
| ---- | ----------- |
| 0    | Magic number (always `0xE9`) |
| 1    | Number of segments |
| 2    | SPI Flash Mode (0 = QIO, 1 = QOUT, 2 = DIO, 3 = DOUT) |
| 3    | High four bits - Flash size (0 = 1MB, 1 = 2MB, 2 = 4MB, 3 = 8MB, 4 = 16MB) |
| 3   | Low four bits - Flash frequency (0 = 40MHz, 1 = 26MHz, 2 = 20MHz, 0xf = 80MHz) |
| 4-7  | Entry point address |

`esptool` overrides the 2nd and 3rd (counted from 0) bytes according to the SPI flash info provided through the command line options (see :ref:`flash-modes`).
These bytes are only overridden if this is a bootloader image (an image written to a correct bootloader offset of {IDF_TARGET_BOOTLOADER_OFFSET}).
In this case, the appended SHA256 digest, which is a cryptographic hash used to verify the integrity of the image, is also updated to reflect the header changes.
Generating images without SHA256 digest can be achieved by running `esptool elf2image` with the `--dont-append-digest` argument.

## Extended File Header

![](./images/firmware_image_ext_header_format.png)

| Byte | Description |
| ---- | ----------- |
| 0    | WP pin when SPI pins set via eFuse (read by ROM bootloader)
| 1-3  | Drive settings for the SPI flash pins (read by ROM bootloader)
| 4-5  | Chip ID (which ESP device is this image for)
| 6    | Minimal chip revision supported by the image (deprecated, use the following field)
| 7-8  | Minimal chip revision supported by the image (in format: major * 100 + minor)
| 9-10 | Maximal chip revision supported by the image (in format: major * 100 + minor)
| 11-14 | Reserved bytes in additional header space, currently unused
| 15 | Hash appended (If 1, SHA256 digest is appended after the checksum)

## Segment

| Byte | Description |
| ---- | ----------- |
| 0-3  | Memory offset |
| 4-7  | Segment size
| 8…n  | Data |


## Footer

The file is padded with zeros until its size is one byte less than a multiple of 16 bytes. A last byte (thus making the file size a multiple of 16) is the checksum of the data of all segments. The checksum is defined as the xor-sum of all bytes and the byte `0xEF`.

If `hash appended` in the extended file header is `0x01`, a SHA256 digest “simple hash” (of the entire image) is appended after the checksum. This digest is separate to secure boot and only used for detecting corruption. The SPI flash info cannot be changed during flashing if hash is appended after the image.

If secure boot is enabled, a signature is also appended (and the simple hash is included in the signed data). This image signature is [Secure Boot V1](https://docs.espressif.com/projects/esp-idf/en/latest/esp32/security/secure-boot-v1.html#image-signing-algorithm) and [Secure Boot V2](https://docs.espressif.com/projects/esp-idf/en/latest/esp32/security/secure-boot-v2.html#signature-block-format) specific.


## Analyzing a Binary Image

To analyze a binary image and get a complete summary of its headers and segments, use the [image-info](https://docs.espressif.com/projects/esptool/en/latest/esp32c3/esptool/basic-commands.html#image-info) command.

