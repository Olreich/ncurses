package ncurses

// struct ldat{};
// struct _win_st{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <curses.h>
import "C"

/* An abstraction of the Mouse. Contains a mouse mask variable and an event variable. */
type Mouse struct {
	mask mMask
	event *C.MEVENT
}

/* Initializes the mouse (must be done after nc.Init). Returns a mouse object, with events enabled from variable mask. */
func InitMouse(mask int) *Mouse {
	mmask := mMask(C.mousemask(C.mmask_t(mask), nil))
	event := C.MEVENT{}
	return &Mouse{mmask, &event}
}
/* Adds events to the mask, making the buttons indicated visible to the input */
func (m *Mouse) AddEvent(mask int) {
	m.mask.addEvents(mask)
}
/* Removes events from the mask, making them invisible */
func (m *Mouse) RemoveEvent(mask int) {
	m.mask.removeEvents(mask)
}
/* Gets a mouse event */
func (m *Mouse) GetEvent() bool {
	if C.getmouse((*C.MEVENT)(m.event)) == 0 { return true }
	return false
}
/* Returns the button state that the event occured on */
func (m *Mouse) State() int {
	return int(m.event.bstate)
}
/* Returns the x coordinate that the mouse event occured on */
func (m *Mouse) X() int {
	return int(m.event.x)
}
/* Returns the y coordinate that the mouse event occured on */
func (m *Mouse) Y() int {
	return int(m.event.y)
}

type mMask C.mmask_t

func (m mMask) addEvents(mask int) {
	m = mMask(C.mousemask(C.mmask_t(mask|int(m)), nil))
}

func (m mMask) removeEvents(mask int) {
	m = mMask(C.mousemask(C.mmask_t(mask &^ int(m)), nil))
}

/* A representation of a window */
type Window C.WINDOW

/* Creates a new window, based on the origin Window (nc.Init returns stdscr, so this should be used for other window). The startx and starty variables begin in the upper-left corner and increase going down and to the right of the screen. These are also relative to the origin window. */
func NewWindow(origin *Window, height, width, startx, starty int) *Window {
	origin.Touch()
	return (*Window)(C.derwin((*C.WINDOW)(origin), C.int(height), C.int(width), C.int(starty), C.int(startx)))
}

/* Removes the window */
func (w *Window) Destroy() {
	w.Clear()
	C.delwin((*C.WINDOW)(w))
}
/* Touches the window, indicating that it needs to be refreshed. It is also required to keep the standard screen from becoming useless once windows become involved (handled by NewWindow) */
func (w *Window) Touch() {
	C.touchwin((*C.WINDOW)(w))
}
/* Refreshes the window */
func (w *Window) Refresh() {
	C.wrefresh((*C.WINDOW)(w))
}
/* Moves cursor to coordinates (x, y) and then prints string s */
func (w *Window) MovePrint(x, y int, s string) {
	C.mvwaddstr((*C.WINDOW)(w), C.int(y), C.int(x), C.CString(s))
}
/* Like MovePrint, but clears the remainder of the line as well */
func (w *Window) MovePrintLine(x, y int, s string) {
	C.mvwaddstr((*C.WINDOW)(w), C.int(y), C.int(x), C.CString(s))
	w.ClearToEOL()
}
/* Prints at the cursor location string s. */
func (w *Window) Print(s string) {
	C.waddstr((*C.WINDOW)(w), C.CString(s))
}
/* Prints at the cursor location string s. Then clears the remainder of the line. */
func (w *Window) PrintLine(s string) {
	w.Print(s)
	w.ClearToEOL()
}
/* Gets a single character from the input stream */
func (w *Window) GetCh() int {
	return int(C.wgetch((*C.WINDOW)(w)))
}
/* Clears the window, erasing any content on it */
func (w *Window) Clear() {
	C.wclear((*C.WINDOW)(w))
}
/* Clears the line the cursor is currently on */
func (w *Window) ClearLine() {
	_, y := w.GetXY()
	w.Move(0, y)
	w.ClearToEOL()
}
/* Clears from the current location of the cursor to the end of the line */
func (w *Window) ClearToEOL() {
	C.wclrtoeol((*C.WINDOW)(w))
}
/* Gets a string from input. Buffering stops when a newline, carriage return, or EOF is encountered. */
func (w *Window) GetStr() (str string) {
	cstr := C.CString(str)
	C.wgetstr((*C.WINDOW)(w), cstr)
	str = C.GoString(cstr)
	return
}
/* Returns the current coordinates of the cursor */
func (w *Window) GetXY() (int, int) {
	return int(C.getcurx((*C.WINDOW)(w))), int(C.getcury((*C.WINDOW)(w)))
	
}
/* Moves the cursor to coordinates (x,y) */
func (w *Window) Move(x, y int) {
	C.wmove((*C.WINDOW)(w), C.int(y), C.int(x))
}
/* Gets the maximum possible coordinates on the current screen */
func (w *Window) GetMaxXY() (x, y int) {
	return int(C.getmaxx((*C.WINDOW)(w))), int(C.getmaxy((*C.WINDOW)(w)))
}
/* Turns an attribute on (one of the A_ constants or a ColorPair). This attribute will remain set until AttrOff is called on that same attribute, or AttrSet is used without the attribute. Attributes can be OR'd together. */
func (w *Window) AttrOn(attr int) {
	C.wattron((*C.WINDOW)(w), C.int(attr))
}
/* Turns an attribute off */
func (w *Window) AttrOff(attr int) {
	C.wattroff((*C.WINDOW)(w), C.int(attr))
}
/* Clears all attributes, then sets an attribute(s). */ 
func (w *Window) AttrSet(attr int) {
	C.wattrset((*C.WINDOW)(w), C.int(attr))
}
/* Surrounds the window by a box, using character v for the vertical parts and character h for the horizontal parts */
func (w *Window) Box(v, h int) {
	C.box((*C.WINDOW)(w), C.chtype(v), C.chtype(h))
}



/* Initializes NCurses library, returns the standard screen. */
func Init() *Window {
	return (*Window)(C.initscr())
}

/* Ends Curses mode, needs to be called before the end of the application, so the terminal returns to normal usage. It is recommended that it be defered in the main curses function. */
func End() {
	C.endwin()
}

/* Sets the Input Mode to raw mode. This means that all characters typed are immediately available to the application without buffering, and reports all control characters directly to the application, skipping the terminal driver altogether. */
func Raw() {
	C.raw()
}
/* Turns off raw input mode.*/
func NoRaw() {
	C.noraw()
}

/* Set the Input Mode to cbreak mode. This means that all characters typed are immediately avaialble to the application though Curses will pass control characters to the terminal before the application. This is the default mode of ncurses. */
func CBreak() {
	C.cbreak()
}

/* Turns off cbreak input mode. */
func NoCBreak() {
	C.nocbreak()
}

/* Enables half-delay mode. This is the same as cbreak mode (see func CBreak), except that there is a time-out delay on the input, and if no input is received, the input will return an Error and move on. Tenths is the delay in tenths of seconds before the program moves on. It must be between 1 and 255. */
func HalfDelay(tenths int) {
	if tenths > 255 { tenths = 255 } else if tenths < 1 { tenths = 1 }
	C.halfdelay(C.int(tenths))
}

/* Echoes incoming characters. */
func Echo() {
	C.echo()
}

/* Turns off echoing, supressing character display that is not explicitly done by the program */
func NoEcho() {
	C.noecho()
}


func DoUpdate() {
	C.doupdate()
}

//-----Input Functions-----//

/* Enables the reading of all keys on a keyboard (that the operating system will interpret) */
func Keypad(win *Window, translate bool) {
	temp := 0
	if translate { temp = 1 }
	C.keypad((*C.WINDOW)(win), C.bool(temp))
}

//-----Color Functions-----//

/* Checks whether the current terminal supports colors, and returns true if it does. */
func HasColors() bool {
	if C.has_colors() > 0 { return true }
	return false
}

/* Checks whether the current terminal supports changing the default color values, and returns true if so */
func CanChangeColors() bool {
	if C.can_change_color() > 0 { return true }
	return false
}

/* Initializes the color engine for ncurses, must be called before InitColor, InitPair, or ColorPair will work */
func StartColor() {
	C.start_color()
}
/* Initializes color number n to red, green, blue values. The color can then be used by InitPair to create Color Pairs*/
func InitColor(n, r, g, b int) {
	C.init_color(C.short(n), C.short(r), C.short(g), C.short(b))
}

/* Initializes a color pair, foreground and background. The color pair number to initialize is n, with the foreground and background colors after that. */
func InitPair(n, fore, back int) {
	C.init_pair(C.short(n), C.short(fore), C.short(back))
}

/* Gets a ColorPair for use in the Attr functions */
func ColorPair(n int) int{
	return int(C.COLOR_PAIR(C.int(n)))
}

const (
	ERR = -0x1;
	OK = 0;
	A_NORMAL = 0;
	A_ALTCHARSET = 0x400000;
	A_BLINK = 0x80000;
	A_BOLD = 0x200000;
	A_DIM = 0x100000;
	A_INVIS = 0x800000;
	A_PROTECT = 0x1000000;
	A_REVERSE = 0x40000;
	A_STANDOUT = 0x10000;
	A_UNDERLINE = 0x20000;
	A_ATTRIBUTES = -0x100;
	A_CHARTEXT = 0xff;
	A_COLOR = 0xff00;
	WA_NORMAL = 0;
	WA_ALTCHARSET = 0x400000;
	WA_BLINK = 0x80000;
	WA_BOLD = 0x200000;
	WA_DIM = 0x100000;
	WA_INVIS = 0x800000;
	WA_LEFT = 0x4000000;
	WA_PROTECT = 0x1000000;
	WA_REVERSE = 0x40000;
	WA_RIGHT = 0x10000000;
	WA_STANDOUT = 0x10000;
	WA_UNDERLINE = 0x20000;
	COLOR_BLACK = 0;
	COLOR_BLUE = 0x4;
	COLOR_GREEN = 0x2;
	COLOR_CYAN = 0x6;
	COLOR_RED = 0x1;
	COLOR_MAGENTA = 0x5;
	COLOR_YELLOW = 0x3;
	COLOR_WHITE = 0x7;
	KEY_BREAK = 0x101;
	KEY_DOWN = 0x102;
	KEY_UP = 0x103;
	KEY_LEFT = 0x104;
	KEY_RIGHT = 0x105;
	KEY_HOME = 0x106;
	KEY_BACKSPACE = 0x107;
	KEY_F0 = 0x108;
	KEY_DL = 0x148;
	KEY_IL = 0x149;
	KEY_DC = 0x14a;
	KEY_IC = 0x14b;
	KEY_EIC = 0x14c;
	KEY_CLEAR = 0x14d;
	KEY_EOS = 0x14e;
	KEY_EOL = 0x14f;
	KEY_SF = 0x150;
	KEY_SR = 0x151;
	KEY_NPAGE = 0x152;
	KEY_PPAGE = 0x153;
	KEY_STAB = 0x154;
	KEY_CTAB = 0x155;
	KEY_CATAB = 0x156;
	KEY_ENTER = 0x157;
	KEY_SRESET = 0x158;
	KEY_RESET = 0x159;
	KEY_PRINT = 0x15a;
	KEY_LL = 0x15b;
	KEY_A1 = 0x15c;
	KEY_A3 = 0x15d;
	KEY_B2 = 0x15e;
	KEY_C1 = 0x15f;
	KEY_C3 = 0x160;
	KEY_BTAB = 0x161;
	KEY_BEG = 0x162;
	KEY_CANCEL = 0x163;
	KEY_CLOSE = 0x164;
	KEY_COMMAND = 0x165;
	KEY_COPY = 0x166;
	KEY_CREATE = 0x167;
	KEY_END = 0x168;
	KEY_EXIT = 0x169;
	KEY_FIND = 0x16a;
	KEY_HELP = 0x16b;
	KEY_MARK = 0x16c;
	KEY_MESSAGE = 0x16d;
	KEY_MOVE = 0x16e;
	KEY_NEXT = 0x16f;
	KEY_OPEN = 0x170;
	KEY_OPTIONS = 0x171;
	KEY_PREVIOUS = 0x172;
	KEY_REDO = 0x173;
	KEY_REFERENCE = 0x174;
	KEY_REFRESH = 0x175;
	KEY_REPLACE = 0x176;
	KEY_RESTART = 0x177;
	KEY_RESUME = 0x178;
	KEY_SAVE = 0x179;
	KEY_SBEG = 0x17a;
	KEY_SCANCEL = 0x17b;
	KEY_SCOMMAND = 0x17c;
	KEY_SCOPY = 0x17d;
	KEY_SCREATE = 0x17e;
	KEY_SDC = 0x17f;
	KEY_SDL = 0x180;
	KEY_SELECT = 0x181;
	KEY_SEND = 0x182;
	KEY_SEOL = 0x183;
	KEY_SEXIT = 0x184;
	KEY_SFIND = 0x185;
	KEY_SHELP = 0x186;
	KEY_SHOME = 0x187;
	KEY_SIC = 0x188;
	KEY_SLEFT = 0x189;
	KEY_SMESSAGE = 0x18a;
	KEY_SMOVE = 0x18b;
	KEY_SNEXT = 0x18c;
	KEY_SOPTIONS = 0x18d;
	KEY_SPREVIOUS = 0x18e;
	KEY_SPRINT = 0x18f;
	KEY_SREDO = 0x190;
	KEY_SREPLACE = 0x191;
	KEY_SRIGHT = 0x192;
	KEY_SRSUME = 0x193;
	KEY_SSAVE = 0x194;
	KEY_SSUSPEND = 0x195;
	KEY_SUNDO = 0x196;
	KEY_SUSPEND = 0x197;
	KEY_UNDO = 0x198;
	KEY_MOUSE = 0631
)

/* Mouse Button Event Constants */
const (
	BUTTON1_RELEASED = 1 << iota
	BUTTON1_PRESSED
	BUTTON1_CLICKED
	BUTTON1_DOUBLE_CLICKED
	BUTTON1_TRIPLE_CLICKED
	_
	BUTTON2_RELEASED
	BUTTON2_PRESSED
	BUTTON2_CLICKED
	BUTTON2_DOUBLE_CLICKED
	BUTTON2_TRIPLE_CLICKED
	_
	BUTTON3_RELEASED
	BUTTON3_PRESSED
	BUTTON3_CLICKED
	BUTTON3_DOUBLE_CLICKED
	BUTTON3_TRIPLE_CLICKED
	_
	BUTTON4_RELEASED
	BUTTON4_PRESSED
	BUTTON4_CLICKED
	BUTTON4_DOUBLE_CLICKED
	BUTTON4_TRIPLE_CLICKED
	_
	BUTTON_CTRL
	BUTTON_SHIFT
	BUTTON_ALT
	REPORT_MOUSE_POSITION
	ALL_MOUSE_EVENTS = REPORT_MOUSE_POSITION - 1
)
