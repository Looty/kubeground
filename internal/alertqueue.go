package internal

import "time"

var MessageQueue []map[string]interface{}

// FormatTime formats the time.Time value to display only hours and minutes.
func FormatTime(t time.Time) string {
	return t.Format("15:04:05") // Example format for hours, minutes, and seconds
}

// AddMapToSlice adds a new map to the slice of maps and returns the updated slice.
func AddMapToSlice(slice []map[string]interface{}, newMap map[string]interface{}) []map[string]interface{} {
	return append(slice, newMap)
}

// RemoveMapFromSlice removes a map from the slice at the specified index and returns the updated slice.
func RemoveMapFromSlice(slice []map[string]interface{}, index int) []map[string]interface{} {
	if index < 0 || index >= len(slice) {
		return slice // Index out of range, return the original slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// Append the new string to the slice and return the updated slice
func AddMessageToQueue(alertType, message string, err string, showSpinner bool) {
	alert := map[string]interface{}{
		"alertType": alertType,
		"message":   message,
		"error":     err,
		"time":      FormatTime(time.Now()),
		"spinner":   showSpinner,
	}
	MessageQueue = AddMapToSlice(MessageQueue, alert)
}

// RemoveMessageFromQueue removes the first map from the MessageQueue that contains the specified message.
func RemoveMessageFromQueue(message string) {
	for i, alert := range MessageQueue {
		if alert["message"] == message {
			MessageQueue = RemoveMapFromSlice(MessageQueue, i)
			break // Remove only the first matching message
		}
	}
}

func RemoveAlertByIndex(index int) {
	if index >= 0 && index < len(MessageQueue) {
		MessageQueue = append(MessageQueue[:index], MessageQueue[index+1:]...)
	}
}

// RemoveAllAlerts clears the MessageQueue slice, effectively removing all alerts.
func RemoveAllAlerts() {
	MessageQueue = []map[string]interface{}{} // Set MessageQueue to an empty slice
}
