package log

var level_codes = map[Level]byte{
	DEBUG:   'D',
	INFO:    'I',
	NOTICE:  'N',
	WARN:    'W',
	ERROR:   'E',
	FATAL:   'F',
	UNKNOWN: 'U',
}
