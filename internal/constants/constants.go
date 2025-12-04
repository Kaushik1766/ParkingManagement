package constants

type ContextKeys string

const User ContextKeys = "user"

const (
	AdminOffice       string = "ADMIN_OFFICE"
	AdminBuilding     string = "ADMIN_BUILDING"
	SlotLayout        string = "000000000000000111111111111111"
	TestBuilding      string = "TEST_BUILDING"
	TestOffice        string = "TEST_OFFICE"
	BillEmailHeader   string = "Your monthly bill generated on %s is here"
	BillEmailTemplate string = `
<!DOCTYPE html>
<html>
<head>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
</head>
<body>

<h2>Monthly Parking Bill</h2>
<p>Dear User,</p>
<p>Here is your parking bill details.</p>

<table>
  <tr>
    <th>Date</th>
    <th>Start Time</th>
    <th>End Time</th>
    <th>Vehicle Type</th>
    <th>Number Plate</th>
    <th>Duration (Hours)</th>
    <th>Cost</th>
  </tr>
  %s
</table>

<h3>Total Amount: %.2f</h3>

<p>Thank you for using our parking service.</p>

</body>
</html>
`
)
