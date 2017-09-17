package deebee

type Manipulation string

const(

	CREATE Manipulation = "CREATE"
	SELECT Manipulation = "SELECT"
	INSERT Manipulation = "INSERT"
	UPDATE Manipulation = "UPDATE"
	DELETE Manipulation = "DELETE"

)

func (m Manipulation) String() string {


	return string(m)
}

