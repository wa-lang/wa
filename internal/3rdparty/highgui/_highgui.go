// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opencv

/*
#include "opencv.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

/*****************************************************************************\
*                         Basic GUI functions                                 *
\*****************************************************************************/

/* this function is used to set some external parameters in case of X Window */
func initSystem(args []string) int {
	argc := C.int(len(args))
	argv := make([]*C.char, len(args))
	for i := 0; i < len(args); i++ {
		argv[i] = C.CString(args[i])
	}
	rv := C.cvInitSystem(argc, (**C.char)(unsafe.Pointer(&argv)))
	return int(rv)
}

func StartWindowThread() int {
	return int(C.cvStartWindowThread())
}

/* wait for key event infinitely (delay<=0) or for "delay" milliseconds */
func WaitKey(delay int) int {
	key := C.cvWaitKey(C.int(delay))
	return int(key)
}

//-----------------------------------------------------------------------------
// Window wrapper for go
//-----------------------------------------------------------------------------

// mouse callback
type MouseFunc interface{}
type MouseFuncA func(event, x, y, flags int)
type MouseFuncB func(event, x, y, flags int, param ...interface{})

// trackbar callback
type TrackbarFunc interface{}
type TrackbarFuncA func(pos int)
type TrackbarFuncB func(pos int, param ...interface{})

// named window
type Window struct {
	name   string
	name_c *C.char
	flags  C.int

	mouseHandle MouseFunc
	param       []interface{}

	trackbarHandle map[string]TrackbarFunc
	trackbarMax    map[string]int
	trackbarVal    map[string]int
	trackbarName   map[string](*C.char)
	trackbarParam  map[string]([]interface{})

	image    *IplImage
	refCount int
}

// window list
var allWindows = make(map[string]*Window, 1000)

//-----------------------------------------------------------------------------
// window: create
//-----------------------------------------------------------------------------

const (
	CV_WINDOW_AUTOSIZE = C.CV_WINDOW_AUTOSIZE
)

/* create window */
func NewWindow(name string, flags ...int) *Window {
	win_flags := C.int(CV_WINDOW_AUTOSIZE)
	if len(flags) > 0 {
		win_flags = C.int(flags[0])
	}

	win := &Window{
		name:   name,
		name_c: C.CString(name),
		flags:  win_flags,

		trackbarHandle: make(map[string]TrackbarFunc, 50),
		trackbarMax:    make(map[string]int, 50),
		trackbarVal:    make(map[string]int, 50),
		trackbarName:   make(map[string](*C.char), 50),
		trackbarParam:  make(map[string]([]interface{}), 50),
	}
	C.cvNamedWindow(win.name_c, win.flags)
	C.GoOpenCV_SetMouseCallback(win.name_c)

	allWindows[win.name] = win
	return win
}

// ---------  YV ---------

/* display image within window (highgui windows remember their content) */
func (win *Window) ShowImage(image *IplImage) {
	win.image = image
	C.cvShowImage(win.name_c, unsafe.Pointer(image))
}

/* resize/move window */
func (win *Window) Resize(width, height int) {
	C.cvResizeWindow(win.name_c, C.int(width), C.int(height))
}
func (win *Window) Move(x, y int) {
	C.cvMoveWindow(win.name_c, C.int(x), C.int(y))
}

/* get native window handle (HWND in case of Win32 and Widget in case of X Window) */
func (win *Window) GetHandle() unsafe.Pointer {
	p := C.cvGetWindowHandle(win.name_c)
	return unsafe.Pointer(p)
}

/* get name of highgui window given its native handle */
func (win *Window) GetWindowName() string {
	return win.name
}

//-----------------------------------------------------------------------------
// window: track bar
//-----------------------------------------------------------------------------

/* create trackbar and display it on top of given window, set callback */
func (win *Window) CreateTrackbar(name string,
	value, count int,
	on_changed TrackbarFunc, param ...interface{}) bool {

	bar_name := C.CString(name)

	switch f := on_changed.(type) {
	case TrackbarFuncA:
		win.trackbarHandle[name] = TrackbarFunc(f)
	case TrackbarFuncB:
		win.trackbarHandle[name] = TrackbarFunc(f)
	case func(pos int):
		win.trackbarHandle[name] = TrackbarFunc(f)
	case func(pos int, param ...interface{}):
		win.trackbarHandle[name] = TrackbarFunc(f)
	default:
		panic("unknow func type!")
	}

	win.trackbarVal[name] = value
	win.trackbarMax[name] = count
	//win.trackbarHandle[name] = on_changed
	win.trackbarName[name] = bar_name

	if len(param) > 0 {
		win.trackbarParam[name] = param
	} else {
		win.trackbarParam[name] = nil
	}

	rv := C.GoOpenCV_CreateTrackbar(bar_name, win.name_c,
		C.int(value), C.int(count))
	return bool(rv != 0)
}
func destroyTrackbar(barName_, winName_ *C.char) {
	C.GoOpenCV_DestroyTrackbar(barName_, winName_)
}

//export goTrackbarCallback
func goTrackbarCallback(barName_, winName_ *C.char, pos C.int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	winName := C.GoString(winName_)
	barName := C.GoString(barName_)

	win, ok := allWindows[winName]
	if !ok {
		return
	}

	trackbarHandle, ok := win.trackbarHandle[barName]
	if !ok {
		return
	}
	if trackbarHandle == nil {
		return
	}

	if fa, ok := trackbarHandle.(func(pos int)); ok {
		fa(int(pos))
	} else if fb, ok := trackbarHandle.(func(pos int, param ...interface{})); ok {
		param := win.trackbarParam[barName]
		if param != nil {
			fb(int(pos), param...)
		} else {
			fb(int(pos))
		}
	}
}

/* retrieve or set trackbar position */
func (win *Window) GetTrackbarPos(name string) (value, max int) {
	rv := C.cvGetTrackbarPos(win.trackbarName[name], win.name_c)
	return int(rv), win.trackbarMax[name]
}
func (win *Window) SetTrackbarPos(name string, pos int) {
	C.cvSetTrackbarPos(win.trackbarName[name], win.name_c, C.int(pos))
}

//-----------------------------------------------------------------------------
// window: mouse callback
//-----------------------------------------------------------------------------

const (
	CV_EVENT_MOUSEMOVE     = int(C.CV_EVENT_MOUSEMOVE)
	CV_EVENT_LBUTTONDOWN   = int(C.CV_EVENT_LBUTTONDOWN)
	CV_EVENT_RBUTTONDOWN   = int(C.CV_EVENT_RBUTTONDOWN)
	CV_EVENT_MBUTTONDOWN   = int(C.CV_EVENT_MBUTTONDOWN)
	CV_EVENT_LBUTTONUP     = int(C.CV_EVENT_LBUTTONUP)
	CV_EVENT_RBUTTONUP     = int(C.CV_EVENT_RBUTTONUP)
	CV_EVENT_MBUTTONUP     = int(C.CV_EVENT_MBUTTONUP)
	CV_EVENT_LBUTTONDBLCLK = int(C.CV_EVENT_LBUTTONDBLCLK)
	CV_EVENT_RBUTTONDBLCLK = int(C.CV_EVENT_RBUTTONDBLCLK)
	CV_EVENT_MBUTTONDBLCLK = int(C.CV_EVENT_MBUTTONDBLCLK)

	CV_EVENT_FLAG_LBUTTON  = int(C.CV_EVENT_FLAG_LBUTTON)
	CV_EVENT_FLAG_RBUTTON  = int(C.CV_EVENT_FLAG_RBUTTON)
	CV_EVENT_FLAG_MBUTTON  = int(C.CV_EVENT_FLAG_MBUTTON)
	CV_EVENT_FLAG_CTRLKEY  = int(C.CV_EVENT_FLAG_CTRLKEY)
	CV_EVENT_FLAG_SHIFTKEY = int(C.CV_EVENT_FLAG_SHIFTKEY)
	CV_EVENT_FLAG_ALTKEY   = int(C.CV_EVENT_FLAG_ALTKEY)
)

/* assign callback for mouse events */
func (win *Window) SetMouseCallback(on_mouse MouseFunc, param ...interface{}) {

	switch f := on_mouse.(type) {
	case MouseFuncA:
		win.mouseHandle = MouseFunc(f)
	case MouseFuncB:
		win.mouseHandle = MouseFunc(f)
	case func(event, x, y, flags int):
		win.mouseHandle = MouseFunc(f)
	case func(event, x, y, flags int, param ...interface{}):
		win.mouseHandle = MouseFunc(f)
	default:
		panic("unknow func type!")
	}

	if len(param) > 0 {
		win.param = param
	} else {
		win.param = nil
	}
}

//export goMouseCallback
func goMouseCallback(name *C.char, event, x, y, flags C.int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	winName := C.GoString(name)
	win, ok := allWindows[winName]

	if !ok {
		return
	}
	if win.mouseHandle == nil {
		return
	}

	if fa, ok := win.mouseHandle.(func(event, x, y, flags int)); ok {
		fa(int(event), int(x), int(y), int(flags))
		return
	}

	if fb, ok := win.mouseHandle.(func(event, x, y, flags int, param ...interface{})); ok {
		if win.param != nil {
			fb(int(event), int(x), int(y), int(flags), win.param...)
		} else {
			fb(int(event), int(x), int(y), int(flags))
		}
		return
	}
}

//-----------------------------------------------------------------------------
// window: destroy
//-----------------------------------------------------------------------------

/* destroy window and all the trackers associated with it */
func (win *Window) Destroy() {
	C.cvDestroyWindow(win.name_c)
	delete(allWindows, win.name)

	for _, bar_name := range win.trackbarName {
		C.GoOpenCV_DestroyTrackbar(bar_name, win.name_c)
		C.free(unsafe.Pointer(bar_name))
	}
	C.free(unsafe.Pointer(win.name_c))
	win.name_c = nil
}

/* destroy window and all the trackers associated with it */
func DestroyAllWindows() {
	for _, win := range allWindows {
		win.Destroy()
	}
}

//-----------------------------------------------------------------------------
// image utils
//-----------------------------------------------------------------------------

const (
	/* 8bit, color or not */
	CV_LOAD_IMAGE_UNCHANGED = int(C.CV_LOAD_IMAGE_UNCHANGED)
	/* 8bit, gray */
	CV_LOAD_IMAGE_GRAYSCALE = int(C.CV_LOAD_IMAGE_GRAYSCALE)
	/* ?, color */
	CV_LOAD_IMAGE_COLOR = int(C.CV_LOAD_IMAGE_COLOR)
	/* any depth, ? */
	CV_LOAD_IMAGE_ANYDEPTH = int(C.CV_LOAD_IMAGE_ANYDEPTH)
	/* ?, any color */
	CV_LOAD_IMAGE_ANYCOLOR = int(C.CV_LOAD_IMAGE_ANYCOLOR)
)

/* load image from file
iscolor can be a combination of above flags where CV_LOAD_IMAGE_UNCHANGED
overrides the other flags
using CV_LOAD_IMAGE_ANYCOLOR alone is equivalent to CV_LOAD_IMAGE_UNCHANGED
unless CV_LOAD_IMAGE_ANYDEPTH is specified images are converted to 8bit
*/

func LoadImage(filename string, iscolor_ ...int) *IplImage {
	iscolor := CV_LOAD_IMAGE_COLOR
	if len(iscolor_) > 0 {
		iscolor = iscolor_[0]
	}
	name_c := C.CString(filename)
	defer C.free(unsafe.Pointer(name_c))

	rv := C.cvLoadImage(name_c, C.int(iscolor))
	return (*IplImage)(rv)
}
func LoadImageM(filename string, iscolor int) *Mat {
	name_c := C.CString(filename)
	defer C.free(unsafe.Pointer(name_c))

	rv := C.cvLoadImageM(name_c, C.int(iscolor))
	return (*Mat)(rv)
}

/* save image to file */
func SaveImage(filename string, image *IplImage) int {
	name_c := C.CString(filename)
	defer C.free(unsafe.Pointer(name_c))
	rv := C.cvSaveImage(name_c, unsafe.Pointer(image))
	return int(rv)
}

const (
	CV_CVTIMG_FLIP    = int(C.CV_CVTIMG_FLIP)
	CV_CVTIMG_SWAP_RB = int(C.CV_CVTIMG_SWAP_RB)
)

/* utility function: convert one image to another with optional vertical flip */
func ConvertImage(src, dst unsafe.Pointer, flags int) {
	C.cvConvertImage(src, dst, C.int(flags))
}

/*****************************************************************************\
*                        Working with Video Files and Cameras                 *
\*****************************************************************************/

/* "black box" capture structure */
type Capture C.CvCapture

/* start capturing frames from video file */
func NewFileCapture(filename string) *Capture {
	filename_c := C.CString(filename)
	defer C.free(unsafe.Pointer(filename_c))
	cap := C.cvCreateFileCapture(filename_c)
	return (*Capture)(cap)
}

const (
	CV_CAP_ANY = int(C.CV_CAP_ANY) // autodetect

	CV_CAP_MIL = int(C.CV_CAP_MIL) // MIL proprietary drivers

	CV_CAP_VFW  = int(C.CV_CAP_VFW) // platform native
	CV_CAP_V4L  = int(C.CV_CAP_V4L)
	CV_CAP_V4L2 = int(C.CV_CAP_V4L2)

	CV_CAP_FIREWARE = int(C.CV_CAP_FIREWARE) // IEEE 1394 drivers
	CV_CAP_FIREWIRE = int(C.CV_CAP_FIREWIRE)
	CV_CAP_IEEE1394 = int(C.CV_CAP_IEEE1394)
	CV_CAP_DC1394   = int(C.CV_CAP_DC1394)
	CV_CAP_CMU1394  = int(C.CV_CAP_CMU1394)

	CV_CAP_STEREO = int(C.CV_CAP_STEREO) // TYZX proprietary drivers
	CV_CAP_TYZX   = int(C.CV_CAP_TYZX)
	CV_TYZX_LEFT  = int(C.CV_TYZX_LEFT)
	CV_TYZX_RIGHT = int(C.CV_TYZX_RIGHT)
	CV_TYZX_COLOR = int(C.CV_TYZX_COLOR)
	CV_TYZX_Z     = int(C.CV_TYZX_Z)

	CV_CAP_QT = int(C.CV_CAP_QT) // QuickTime

	CV_CAP_UNICAP = int(C.CV_CAP_UNICAP) // Unicap drivers

	CV_CAP_DSHOW = int(C.CV_CAP_DSHOW) // DirectShow (via videoInput)
)

/* start capturing frames from camera: index = camera_index + domain_offset (CV_CAP_*) */
func NewCameraCapture(index int) *Capture {
	cap := C.cvCreateCameraCapture(C.int(index))
	return (*Capture)(cap)
}

/* grab a frame, return 1 on success, 0 on fail.
this function is thought to be fast               */
func (capture *Capture) GrabFrame() bool {
	rv := C.cvGrabFrame((*C.CvCapture)(capture))
	return (rv != C.int(0))
}

/* get the frame grabbed with cvGrabFrame(..)
This function may apply some frame processing like
frame decompression, flipping etc.
!!!DO NOT RELEASE or MODIFY the retrieved frame!!! */
func (capture *Capture) RetrieveFrame() *IplImage {
	rv := C.cvRetrieveFrame((*C.CvCapture)(capture))
	return (*IplImage)(rv)
}

/* Just a combination of cvGrabFrame and cvRetrieveFrame
   !!!DO NOT RELEASE or MODIFY the retrieved frame!!!      */
func (capture *Capture) QueryFrame() *IplImage {
	rv := C.cvQueryFrame((*C.CvCapture)(capture))
	return (*IplImage)(rv)
}

/* stop capturing/reading and free resources */
func (capture *Capture) Release() {
	cap_c := (*C.CvCapture)(capture)
	C.cvReleaseCapture(&cap_c)
}

const (
	CV_CAP_PROP_POS_MSEC      = int(C.CV_CAP_PROP_POS_MSEC)
	CV_CAP_PROP_POS_FRAMES    = int(C.CV_CAP_PROP_POS_FRAMES)
	CV_CAP_PROP_POS_AVI_RATIO = int(C.CV_CAP_PROP_POS_AVI_RATIO)
	CV_CAP_PROP_FRAME_WIDTH   = int(C.CV_CAP_PROP_FRAME_WIDTH)
	CV_CAP_PROP_FRAME_HEIGHT  = int(C.CV_CAP_PROP_FRAME_HEIGHT)
	CV_CAP_PROP_FPS           = int(C.CV_CAP_PROP_FPS)
	CV_CAP_PROP_FOURCC        = int(C.CV_CAP_PROP_FOURCC)
	CV_CAP_PROP_FRAME_COUNT   = int(C.CV_CAP_PROP_FRAME_COUNT)
	CV_CAP_PROP_FORMAT        = int(C.CV_CAP_PROP_FORMAT)
	CV_CAP_PROP_MODE          = int(C.CV_CAP_PROP_MODE)
	CV_CAP_PROP_BRIGHTNESS    = int(C.CV_CAP_PROP_BRIGHTNESS)
	CV_CAP_PROP_CONTRAST      = int(C.CV_CAP_PROP_CONTRAST)
	CV_CAP_PROP_SATURATION    = int(C.CV_CAP_PROP_SATURATION)
	CV_CAP_PROP_HUE           = int(C.CV_CAP_PROP_HUE)
	CV_CAP_PROP_GAIN          = int(C.CV_CAP_PROP_GAIN)
	CV_CAP_PROP_EXPOSURE      = int(C.CV_CAP_PROP_EXPOSURE)
	CV_CAP_PROP_CONVERT_RGB   = int(C.CV_CAP_PROP_CONVERT_RGB)
	// CV_CAP_PROP_WHITE_BALANCE = int(C.CV_CAP_PROP_WHITE_BALANCE)
	CV_CAP_PROP_RECTIFICATION = int(C.CV_CAP_PROP_RECTIFICATION)
)

/* retrieve or set capture properties */
func (capture *Capture) GetProperty(property_id int) float64 {
	rv := C.cvGetCaptureProperty((*C.CvCapture)(capture),
		C.int(property_id),
	)
	return float64(rv)
}
func (capture *Capture) SetProperty(property_id int, value float64) int {
	rv := C.cvSetCaptureProperty((*C.CvCapture)(capture),
		C.int(property_id), C.double(value),
	)
	return int(rv)
}

//-----------------------------------------------------------------------------
// VideoWriter
//-----------------------------------------------------------------------------

/* "black box" video file writer structure */
type VideoWriter C.CvVideoWriter

// Prototype for CV_FOURCC so that swig can generate wrapper without mixing up the define
func FOURCC(c1, c2, c3, c4 int8) int32 {
	rv := C.GoOpenCV_FOURCC_(C.int(c1), C.int(c2), C.int(c3), C.int(c4))
	return int32(rv)
}

const (
	/* Open Codec Selection Dialog (Windows only) */
	CV_FOURCC_PROMPT = int(C.CV_FOURCC_PROMPT)
	/* Use default codec for specified filename (Linux only) */
	// CV_FOURCC_DEFAULT = int(C.CV_FOURCC_DEFAULT)
)

/* initialize video file writer */
func NewVideoWriter(filename string,
	fourcc int32, fps float32,
	frame_width, frame_height,
	is_color int) *VideoWriter {

	size := C.cvSize(C.int(frame_width), C.int(frame_height))
	filename_c := C.CString(filename)
	defer C.free(unsafe.Pointer(filename_c))

	rv := C.cvCreateVideoWriter(filename_c,
		C.int(fourcc), C.double(fps), size, C.int(is_color),
	)
	return (*VideoWriter)(rv)
}

/* write frame to video file */
func (writer *VideoWriter) WriteFrame(image *IplImage) int {
	rv := C.cvWriteFrame((*C.CvVideoWriter)(writer), (*C.IplImage)(image))
	return int(rv)
}

/* close video file writer */
func (writer *VideoWriter) Release() {
	writer_c := (*C.CvVideoWriter)(writer)
	C.cvReleaseVideoWriter(&writer_c)
}

/*****************************************************************************\
*                                 --- END ---                                 *
\*****************************************************************************/
