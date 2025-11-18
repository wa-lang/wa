# Boot Mode Selection

*https://docs.espressif.com/projects/esptool/en/latest/esp32c3/advanced-topics/boot-mode-selection.html*

This guide explains how to select the boot mode correctly and describes the boot log messages of ESP32-C3.

> Warning: The ESP32-C3 has a 45k ohm internal pull-up/pull-down resistor at GPIO9 (and other pins). If you want to connect a switch button to enter the boot mode, this has to be a strong pull-down. For example a 10k resistor to GND.

Information about ESP32-C3 strapping pins can also be found in the [ESP32-C3 Datasheet](https://www.espressif.com/sites/default/files/documentation/esp32-c3_datasheet_en.pdf), section “Strapping Pins”.

On many development boards with built-in USB/Serial, `esptool` can automatically reset the board into bootloader mode. For other configurations or custom hardware, you will need to check the orientation of some “strapping pins” to get the correct boot mode:

## Select Bootloader Mode

### GPIO9

The ESP32-C3 will enter the serial bootloader when GPIO9 is held low on reset. Otherwise it will run the program in flash.

| GPIO9 Input | Mode |
| ----------- | ---- |
| Low/GND     | ROM serial bootloader for esptool
| High/VCC    | Normal execution mode

GPIO9 has an internal pullup resistor, so if it is left unconnected then it will pull high.

Many boards use a button marked “Flash” (or “BOOT” on some Espressif development boards) that pulls GPIO9 low when pressed.

### GPIO8

GPIO8 must also be driven High, in order to enter the serial bootloader reliably. The strapping combination of GPIO8 = 0 and GPIO9 = 0 is invalid and will trigger unexpected behavior.

In normal boot mode (GPIO9 high), GPIO8 is ignored.

### Other Pins

As well as the above mentioned pins, other ones influence the serial bootloader, please consult the [ESP32-C3 Datasheet](https://www.espressif.com/sites/default/files/documentation/esp32-c3_datasheet_en.pdf), section “Strapping Pins”.

## Automatic Bootloader

`esptool` resets ESP32-C3 automatically by asserting `DTR` and `RTS` control lines of the USB to serial converter chip, i.e., FTDI, CP210x, or CH340x. The `DTR` and `RTS` control lines are in turn connected to `GPIO9` and `EN` (`CHIP_PU`) pins of ESP32-C3, thus changes in the voltage levels of `DTR` and `RTS` will boot the ESP32-C3 into Firmware Download mode.

> Note: When developing esptool, keep in mind DTR and RTS are active low signals, i.e., True = pin @ 0V, False = pin @ VCC.

As an example of auto-reset curcuitry implementation, check the [schematic](https://dl.espressif.com/dl/schematics/esp32_devkitc_v4-sch-20180607a.pdf) of the ESP32 DevKitC development board:

- The Micro USB 5V & USB-UART section shows the `DTR` and `RTS` control lines of the USB to serial converter chip connected to `GPIO9` and `EN` pins of the ESP module.
- Some OS and/or drivers may activate `RTS` and or `DTR` automatically when opening the serial port (true only for some serial terminal programs, not `esptool`), pulling them low together and holding the ESP in reset. If `RTS` is wired directly to `EN` then RTS/CTS “hardware flow control” needs to be disabled in the serial program to avoid this. An additional circuitry is implemented in order to avoid this problem - if both `RTS` and `DTR` are asserted together, this doesn’t reset the chip. The schematic shows this specific circuit with two transistors and its truth table.
- If this circuitry is implemented (all Espressif boards have it), adding a capacitor between the `EN` pin and `GND` (in the 1uF-10uF range) is necessary for the reset circuitry to work reliably. This is shown in the ESP32 Module section of the schematic.
- The Switch Button section shows buttons needed for [manually switching to bootloader](https://docs.espressif.com/projects/esptool/en/latest/esp32c3/advanced-topics/boot-mode-selection.html#manual-bootloader).

Make the following connections for `esptool` to automatically enter the bootloader of an ESP32-C3 chip:

| ESP Pin | Serial Pin |
| ------- | ---------- |
| EN | RTS
| GPIO9 | DTR

In Linux serial ports by default will assert RTS when nothing is attached to them. This can hold the ESP32-C3 in a reset loop which may cause some serial adapters to subsequently reset loop. This functionality can be disabled by disabling `HUPCL` (ie `sudo stty -F /dev/ttyUSB0 -hupcl`).

(Some third party ESP32-C3 development boards use an automatic reset circuit for `EN` & `GPIO9` pins, but don’t add a capacitor on the EN pin. This results in unreliable automatic reset, especially on Windows. Adding a 1uF (or higher) value capacitor between `EN` pin and `GND` may make automatic reset more reliable.)

In general, you should have no problems with the official Espressif development boards. However, `esptool` is not able to reset your hardware automatically in the following cases:

- Your hardware does not have the `DTR` and `RTS` lines connected to `GPIO9` and `EN` (`CHIP_PU`)
- The `DTR` and `RTS` lines are configured differently
- There are no such serial control lines at all

## Manual Bootloader

Depending on the kind of hardware you have, it may also be possible to manually put your ESP32-C3 board into Firmware Download mode (reset).

- For development boards produced by Espressif, this information can be found in the respective getting started guides or user guides. For example, to manually reset a development board, hold down the Boot button (`GPIO9`) and press the EN button (`EN` (`CHIP_PU`)).
- For other types of hardware, try pulling `GPIO9` down.

> Note: If esptool is able to reset the chip but for some reason the chip is not entering into bootloader mode then hold down the Boot button (or pull down `GPIO9`) while you start esptool and keep it down during reset.

## Boot Log

### Boot Mode Message

After reset, the second line printed by the ESP32-C3 ROM (at 115200bps) is a reset & boot mode message:

```
ets Jun  8 2016 00:22:57
rst:0x1 (POWERON_RESET),boot:0x3 (DOWNLOAD_BOOT(UART0/UART1/SDIO_REI_REO_V2))
```

`rst:0xNN (REASON)` is an enumerated value (and description) of the reason for the reset. A mapping between the hex value and each reason can be found in the [ESP-IDF source under RESET_REASON enum](https://github.com/espressif/esp-idf/blob/release/v5.2/components/esp_rom/include/esp32c3/rom/rtc.h). The value can be read in ESP32-C3 code via the [get_reset_reason() ROM function](https://github.com/espressif/esp-idf/blob/release/v5.2/components/esp_rom/include/esp32c3/rom/rtc.h).

`boot:0xNN (DESCRIPTION)` is the hex value of the strapping pins, as represented in the [GPIO_STRAP register](https://github.com/espressif/esp-idf/blob/release/v5.2/components/soc/esp32c3/include/soc/gpio_reg.h).

The individual bit values are as follows:

- `0x04` - GPIO8
- `0x08` - GPIO9

If the pin was high on reset, the bit value will be set. If it was low on reset, the bit will be cleared.

A number of boot mode strings can be shown depending on which bits are set:

- `DOWNLOAD_BOOT(UART0/UART1/SDIO_REI_REO_V2)` or `DOWNLOAD(USB/UART0)` - ESP32-C3 is in download flashing mode (suitable for esptool)
- `SPI_FAST_FLASH_BOOT` - This is the normal SPI flash boot mode.
- Other modes (including `SPI_FLASH_BOOT`, `SDIO_REI_FEO_V1_BOOT`, `ATE_BOOT`) may be shown here. This indicates an unsupported boot mode has been selected. Consult the strapping pins shown above (in most cases, one of these modes is selected if GPIO8 has been pulled high when GPIO9 is low).

### Later Boot Messages

Later output from the ROM bootloader depends on the strapping pins and the boot mode. Some common output includes:

#### Early Flash Read Error

```
Invalid header <value at 0x0>
```

This fatal error indicates that the bootloader tried to read the software bootloader header at address 0x0 but failed to read valid data. Possible reasons for this include:

- There isn’t actually a bootloader at offset 0x0 (maybe the bootloader was flashed to the wrong offset by mistake, or the flash has been erased and no bootloader has been flashed yet.)
- Physical problem with the connection to the flash chip, or flash chip power.
- Flash encryption is enabled but the bootloader is plaintext. Alternatively, flash encryption is disabled but the bootloader is encrypted ciphertext.

#### Software Bootloader Header Info

```
SPIWP:0xee
mode:DIO, clock div:1
```

This is normal boot output based on a combination of eFuse values and information read from the bootloader header at flash offset 0x0:

- `SPIWP:0xNN` indicates a custom `WP` pin value, which is stored in the bootloader header. This pin value is only used if SPI flash pins have been remapped via eFuse (as shown in the `configsip` value). All custom pin values but WP are encoded in the configsip byte loaded from eFuse, and WP is supplied in the bootloader header.
- `mode: AAA, clock div: N`. SPI flash access mode. Read from the bootloader header, correspond to the `--flash-mode` and `--flash-freq` arguments supplied to `esptool write-flash` or `esptool elf2image`.
- `mode` can be DIO, DOUT, QIO, or QOUT. QIO and QOUT are not supported here, to boot in a Quad I/O mode the ROM bootloader should load the software bootloader in a Dual I/O mode and then the ESP-IDF software bootloader enables Quad I/O based on the detected flash chip mode.
- `clock div: N` is the SPI flash clock frequency divider. This is an integer clock divider value from an 80MHz APB clock, based on the supplied `--flash-freq` argument (ie 80MHz=1, 40MHz=2, etc). The ROM bootloader actually loads the software bootloader at a lower frequency than the `--flash-freq` value. The initial APB clock frequency is equal to the crystal frequency, so with a 40MHz crystal the SPI clock used to load the software bootloader will be half the configured value (40MHz/2=20MHz). When the software bootloader starts it sets the APB clock to 80MHz causing the SPI clock frequency to match the value set when flashing.

#### Software Bootloader Load Segments

```
load:0x3fff0008,len:8
load:0x3fff0010,len:3680
load:0x40078000,len:8364
load:0x40080000,len:252
entry 0x40080034
```

These entries are printed as the ROM bootloader loads each segment in the software bootloader image. The load address and length of each segment is printed.

You can compare these values to the software bootloader image by running `esptool --chip esp32c3 image-info /path/to/bootloader.bin` to dump image info including a summary of each segment. Corresponding details will also be found in the bootloader ELF file headers.

If there is a problem with the SPI flash chip addressing mode, the values printed by the bootloader here may be corrupted.

The final line shows the entry point address of the software bootloader, where the ROM bootloader will call as it hands over control.

