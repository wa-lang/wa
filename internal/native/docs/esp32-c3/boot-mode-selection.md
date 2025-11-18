# 启动模式选择

本指南介绍如何正确选择启动模式，并描述 ESP32-C3 的启动日志消息。

> 警告：ESP32-C3 的 GPIO9（和其他引脚）具有 45k 欧姆的内部上拉/下拉电阻。如果要连接一个开关按钮来进入启动模式，这必须是一个强下拉。例如一个 10k 电阻接地。

关于 ESP32-C3 绑定位的信息也可以在 [ESP32-C3 数据手册](https://www.espressif.com/sites/default/files/documentation/esp32-c3_datasheet_en.pdf) 的“绑定位”部分找到。

在许多带有内置 USB/Serial 的开发板上，`esptool` 可以自动将板重置为引导加载程序模式。对于其他配置或定制硬件，您需要检查一些“绑定位”的状态以获得正确的启动模式：

## 选择引导加载程序模式

### GPIO9

当 GPIO9 在复位时保持低电平时，ESP32-C3 将进入串行引导加载程序。否则，它将运行闪存中的程序。

| GPIO9 输入 | 模式 |
| ----------- | ---- |
| 低电平/GND | 用于 esptool 的 ROM 串行引导加载程序 |
| 高电平/VCC | 正常执行模式 |

GPIO9 有一个内部上拉电阻，所以如果它未连接，它将被拉为高电平。

许多开发板使用一个标有“Flash”（或在一些 Espressif 开发板上标为“BOOT”）的按钮，按下时将 GPIO9 拉低。

### GPIO8

为了可靠地进入串行引导加载程序，GPIO8 也必须被驱动为高电平。GPIO8 = 0 和 GPIO9 = 0 的绑定组合是无效的，会触发意外行为。

在正常启动模式（GPIO9 高电平）下，GPIO8 被忽略。

### 其他引脚

除了上述引脚外，其他引脚也会影响串行引导加载程序，请参考 [ESP32-C3 数据手册](https://www.espressif.com/sites/default/files/documentation/esp32-c3_datasheet_en.pdf) 的“绑定位”部分。

## 自动引导加载程序

`esptool` 通过断言 USB 到串行转换器芯片（即 FTDI、CP210x 或 CH340x）的 `DTR` 和 `RTS` 控制线来自动重置 ESP32-C3。`DTR` 和 `RTS` 控制线又连接到 ESP32-C3 的 `GPIO9` 和 `EN`（`CHIP_PU`）引脚，因此 `DTR` 和 `RTS` 的电压电平变化将使 ESP32-C3 进入固件下载模式。

> 注意：在开发 esptool 时，请记住 DTR 和 RTS 是低电平有效信号，即 True = 引脚 @ 0V，False = 引脚 @ VCC。

作为自动复位电路实现的示例，请查看 ESP32 DevKitC 开发板的 [原理图](https://dl.espressif.com/dl/schematics/esp32_devkitc_v4-sch-20180607a.pdf)：

- Micro USB 5V 和 USB-UART 部分显示 USB 到串行转换器芯片的 `DTR` 和 `RTS` 控制线连接到 ESP 模块的 `GPIO9` 和 `EN` 引脚。
- 一些操作系统和/或驱动程序在打开串行端口时可能会自动激活 `RTS` 和/或 `DTR`（仅适用于某些串行终端程序，不适用于 `esptool`），将它们一起拉低并使 ESP 保持复位状态。如果 `RTS` 直接连接到 `EN`，则需要在串行程序中禁用 RTS/CTS“硬件流控制”以避免这种情况。为了避免这个问题，实现了一个额外的电路——如果 `RTS` 和 `DTR` 同时被断言，这不会重置芯片。原理图显示了这个带有两个晶体管的特定电路及其真值表。
- 如果实现了此电路（所有 Espressif 开发板都有），则需要在 `EN` 引脚和 `GND` 之间添加一个电容器（在 1uF-10uF 范围内），以使复位电路可靠工作。这在原理图的 ESP32 模块部分显示。
- 开关按钮部分显示了用于 [手动切换到引导加载程序](https://docs.espressif.com/projects/esptool/en/latest/esp32c3/advanced-topics/boot-mode-selection.html#manual-bootloader) 所需的按钮。

为使 `esptool` 自动进入 ESP32-C3 芯片的引导加载程序，请进行以下连接：

| ESP 引脚 | 串行引脚 |
| ------- | ---------- |
| EN | RTS |
| GPIO9 | DTR |

在 Linux 中，默认情况下，当没有任何设备连接到串行端口时，串行端口会断言 RTS。这可能会使 ESP32-C3 处于复位循环中，这可能导致一些串行适配器随后进入复位循环。可以通过禁用 `HUPCL` 来禁用此功能（即 `sudo stty -F /dev/ttyUSB0 -hupcl`）。

（一些第三方 ESP32-C3 开发板使用 `EN` 和 `GPIO9` 引脚的自动复位电路，但没有在 EN 引脚上添加电容器。这会导致自动复位不可靠，尤其是在 Windows 上。在 `EN` 引脚和 `GND` 之间添加一个 1uF（或更高）值的电容器可能会使自动复位更可靠。）

一般来说，官方 Espressif 开发板应该没有问题。但是，在以下情况下，`esptool` 无法自动重置您的硬件：

- 您的硬件没有将 `DTR` 和 `RTS` 线连接到 `GPIO9` 和 `EN`（`CHIP_PU`）
- `DTR` 和 `RTS` 线配置不同
- 根本没有此类串行控制线

## 手动引导加载程序

根据您拥有的硬件类型，也可以手动将 ESP32-C3 开发板置于固件下载模式（复位）。

- 对于 Espressif 生产的开发板，此信息可以在相应的入门指南或用户指南中找到。例如，要手动重置开发板，请按住 Boot 按钮（`GPIO9`）并按下 EN 按钮（`EN`（`CHIP_PU`））。
- 对于其他类型的硬件，请尝试将 `GPIO9` 拉低。

> 注意：如果 esptool 能够重置芯片，但由于某种原因芯片没有进入引导加载程序模式，请在启动 esptool 时按住 Boot 按钮（或下拉 `GPIO9`），并在复位期间保持按住。

## 启动日志

### 启动模式消息

复位后，ESP32-C3 ROM 打印的第二行（波特率 115200）是复位和启动模式消息：

```
ets Jun  8 2016 00:22:57
rst:0x1 (POWERON_RESET),boot:0x3 (DOWNLOAD_BOOT(UART0/UART1/SDIO_REI_REO_V2))
```

`rst:0xNN (REASON)` 是复位原因的枚举值（和描述）。十六进制值和每个原因之间的映射可以在 [ESP-IDF 源代码中的 RESET_REASON 枚举](https://github.com/espressif/esp-idf/blob/release/v5.2/components/esp_rom/include/esp32c3/rom/rtc.h) 中找到。可以通过 [get_reset_reason() ROM 函数](https://github.com/espressif/esp-idf/blob/release/v5.2/components/esp_rom/include/esp32c3/rom/rtc.h) 在 ESP32-C3 代码中读取该值。

`boot:0xNN (DESCRIPTION)` 是绑定位的十六进制值，如 [GPIO_STRAP 寄存器](https://github.com/espressif/esp-idf/blob/release/v5.2/components/soc/esp32c3/include/soc/gpio_reg.h) 中所表示。

各个位值如下：

- `0x04` - GPIO8
- `0x08` - GPIO9

如果引脚在复位时为高电平，则位值将被设置。如果在复位时为低电平，则该位将被清除。

根据设置的位，可以显示多个启动模式字符串：

- `DOWNLOAD_BOOT(UART0/UART1/SDIO_REI_REO_V2)` 或 `DOWNLOAD(USB/UART0)` - ESP32-C3 处于下载闪存模式（适用于 esptool）
- `SPI_FAST_FLASH_BOOT` - 这是正常的 SPI 闪存启动模式。
- 可能会显示其他模式（包括 `SPI_FLASH_BOOT`、`SDIO_REI_FEO_V1_BOOT`、`ATE_BOOT`）。这表明已选择不支持的启动模式。请参考上面显示的绑定位（在大多数情况下，如果 GPIO9 为低电平时 GPIO8 被拉高，则会选择这些模式之一）。

### 后续启动消息

ROM 引导加载程序的后续输出取决于绑定位和启动模式。一些常见的输出包括：

#### 早期闪存读取错误

```
Invalid header <value at 0x0>
```

此致命错误表明引导加载程序尝试读取地址 0x0 处的软件引导加载程序头，但未能读取有效数据。可能的原因包括：

- 在偏移量 0x0 处实际上没有引导加载程序（可能是引导加载程序被错误地闪存到错误的偏移量，或者闪存已被擦除且尚未闪存任何引导加载程序。）
- 与闪存芯片的连接或闪存芯片电源存在物理问题。
- 启用了闪存加密，但引导加载程序是明文。或者，禁用了闪存加密，但引导加载程序是加密的密文。

#### 软件引导加载程序头信息

```
SPIWP:0xee
mode:DIO, clock div:1
```

这是基于 eFuse 值和从闪存偏移量 0x0 处的引导加载程序头读取的信息的正常启动输出：

- `SPIWP:0xNN` 表示自定义的 `WP` 引脚值，存储在引导加载程序头中。仅当通过 eFuse 重新映射 SPI 闪存引脚时（如 `configsip` 值所示），才会使用此引脚值。除 WP 之外的所有自定义引脚值都编码在从 eFuse 加载的 configsip 字节中，WP 在引导加载程序头中提供。
- `mode: AAA, clock div: N`。SPI 闪存访问模式。从引导加载程序头读取，对应于提供给 `esptool write-flash` 或 `esptool elf2image` 的 `--flash-mode` 和 `--flash-freq` 参数。
- `mode` 可以是 DIO、DOUT、QIO 或 QOUT。这里不支持 QIO 和 QOUT，要以 Quad I/O 模式启动，ROM 引导加载程序应以来自 I/O 模式加载软件引导加载程序，然后 ESP-IDF 软件引导加载程序根据检测到的闪存芯片模式启用 Quad I/O。
- `clock div: N` 是 SPI 闪存时钟频率分频器。这是来自 80MHz APB 时钟的整数时钟分频器值，基于提供的 `--flash-freq` 参数（即 80MHz=1，40MHz=2 等）。ROM 引导加载程序实际上以低于 `--flash-freq` 值的频率加载软件引导加载程序。初始 APB 时钟频率等于晶体频率，因此对于 40MHz 晶体，用于加载软件引导加载程序的 SPI 时钟将是配置值的一半（40MHz/2=20MHz）。当软件引导加载程序启动时，它将 APB 时钟设置为 80MHz，使 SPI 时钟频率与闪存时设置的值匹配。

#### 软件引导加载程序加载段

```
load:0x3fff0008,len:8
load:0x3fff0010,len:3680
load:0x40078000,len:8364
load:0x40080000,len:252
entry 0x40080034
```

这些条目在 ROM 引导加载程序加载软件引导加载程序映像中的每个段时打印。打印每个段的加载地址和长度。

您可以通过运行 `esptool --chip esp32c3 image-info /path/to/bootloader.bin` 来将这些值与软件引导加载程序映像进行比较，以转储映像信息，包括每个段的摘要。相应的详细信息也可以在引导加载程序 ELF 文件头中找到。

如果 SPI 闪存芯片寻址模式有问题，引导加载程序在此处打印的值可能会损坏。

最后一行显示软件引导加载程序的入口点地址，ROM 引导加载程序在移交控制权时将调用该地址。

