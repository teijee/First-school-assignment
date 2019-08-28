package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type bestellingen_s struct {
	bestelnummer            int
	status                  string
	besteldatum             string
	afbetaling_doorlooptijd int
	afbetaling_maandbedrag  float64
	klantnummer             int
	verkoper                int
}

type module struct {
	modulenaam string
	stukprijs  float64
}

type configuratie struct {
	bestelnummer int
	modulenaam   string
	verkoopprijs float64
}

type Klant struct {
	klantnummer           int
	naam                  string
	voornaam              string
	postcode              string
	huisnummer            int
	huisnummer_toevoeging string
	geboortedatum         string
	geslacht              string
	bloedgroep            string
	rhesusfactor          string
	beroepsrisicofactor   float64
	inkomen               int
	kredietregistratie    string
	opleiding             string
	opmerking             string
}

type medewerker struct {
	naam            string
	datum_in_dienst string

	//misschien deze verwijderen als oude code begint te miepen
	medewerkernummer int
}

func main() {

	//scanner defineren
	scanner := bufio.NewScanner(os.Stdin)

	//inlog van de gebruiker
	var gebruikersnaam string
	var datum_in_dienst string

	//inlog "module"
	for inlog := true; inlog == true; {
		fmt.Println("Welkom bij het administratie systeem van Vita Intellect")
		fmt.Println("Om in te loggen hebben wij uw voornaam en datum van indiensttreding nodig")
		fmt.Println("Uw voornaam: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		gebruikersnaam = scanner.Text()
		fmt.Println("Datum van indiensttreding (in YYYY-MM-DD): ")
		scanner.Scan()
		datum_in_dienst = scanner.Text()
		inlog = CheckInlog(gebruikersnaam, datum_in_dienst)
	}
	//keuze menu
	for {
		fmt.Println("U heeft de volgende keuze mogelijkehden: ")
		fmt.Println("Toets 1: Om alle klanten in te zien")
		fmt.Println("Toets 2: Om een Bestelling voor een nieuwe klant toe voegen")
		fmt.Println("Toets 3: Om een Bestelling voor een bestaande klant toe te voegen")
		fmt.Println("Toest 4: Om alle openstaande bestellingen in te zien")
		fmt.Println("Toets 5: Om programma af te sluiten")
		scanner.Scan()
		menukeuze := scanner.Text()
	loop2:
		for {
			switch menukeuze {
			case "1":
				AlleKlanten()
				fmt.Println("Type X om terug te keren naar het hoofdmenu en N om het overzicht opnieuw te laden")
				scanner.Scan()
				menukeuze := scanner.Text()

				if menukeuze == "N" {
					AllOpenOrders()
				} else if menukeuze == "X" {
					break loop2
				}

			case "2":
				nieuweKlanttoevoegen(gebruikersnaam, datum_in_dienst)
				break loop2
			case "3":
				bestaandeklant(gebruikersnaam, datum_in_dienst)
				break loop2
			case "4":
				AllOpenOrders()
				fmt.Println("Type X om terug te keren naar het hoofdmenu en N om het overzicht opnieuw te laden")
				scanner.Scan()
				menukeuze := scanner.Text()

				if menukeuze == "N" {
					AllOpenOrders()
				} else if menukeuze == "X" {
					break loop2
				}
			case "5":
				os.Exit(1)

			}
		}
	}
}

func CheckInlog(gebruikersnaam string, wachtwoord string) bool {

	//vebinding met database opzetten
	db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//de query om de gegvens van de medewerker op te vragen
	query := "SELECT naam, datum_in_dienst FROM vitaintellectdb.medewerker WHERE naam = ? and functie LIKE '%verkoper%'"

	dbresultaat, err := db.Query(query, gebruikersnaam)
	if err != nil {
		panic(err)
	}

	var verkoper medewerker
	if dbresultaat.Next() {
		err := dbresultaat.Scan(&verkoper.naam, &verkoper.datum_in_dienst)
		if err != nil {
			panic(err)
		}
	}

	var trueorfalse bool

	if gebruikersnaam == "" && wachtwoord == "" {
		trueorfalse = true
		fmt.Printf("Inlog gegevens onjuist ,probeer het opnieuw\n\n")
	} else if verkoper.naam == gebruikersnaam && verkoper.datum_in_dienst == wachtwoord {
		trueorfalse = false
		fmt.Println("Inlog succesvol")
		fmt.Println("Welkom ", verkoper.naam)
	} else {
		trueorfalse = true
		fmt.Printf("Inlog gegevens onjuist, probeer het opnieuw \n\n")
	}

	dbresultaat.Close()
	return trueorfalse
}

func nieuweKlanttoevoegen(gebruikersnaam, datum_in_dienst string) {
loop:
	for {
		db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")

		// if there is an error opening the connection, handle it
		if err != nil {
			panic(err.Error())
		}

		defer db.Close()

		fmt.Println("Je wilt een nieuwe bestelling toevoegen van een nieuwe klant.")
		fmt.Println("Deze stappen zullen wij met je doorlopen.")

		var nieuweklant Klant

		//naam van de klant vragen
		fmt.Println("Wat is de achternaam van de klant?")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		nieuweklant.naam = scanner.Text()

		//Achternaam van de klant toevoegen
		fmt.Println("Wat is de voornaam van de klant?")
		scanner.Scan()
		nieuweklant.voornaam = scanner.Text()

		//Postcode van de klant vragen
		fmt.Println("Wat is de postcode van de klant? (1111XX)")
		scanner.Scan()
		nieuweklant.postcode = scanner.Text()

		//Huisnummer van de klant vragen
		fmt.Println("Wat is het huisnummer van de klant?")
		scanner.Scan()
		huisnummer := scanner.Text()
		nieuweklant.huisnummer, err = strconv.Atoi(huisnummer)

		//huisnummer toevoeging van de klant
		fmt.Println("Is er een huisnummer toevoeging? Type Y voor ja en N voor nee")
		scanner.Scan()
		welofgeenhuisnummertoevoeging := scanner.Text()
		nieuweklant.huisnummer_toevoeging = ""
		if welofgeenhuisnummertoevoeging == "Y" {
			fmt.Println("Wat is de huisnummertoevoeging van de klant?")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			nieuweklant.huisnummer_toevoeging = scanner.Text()
		}

		//Geboortedatum invoeren
		fmt.Println("Wat is de geboortedatum van de klant? YYYY-MM-DD")
		scanner.Scan()
		nieuweklant.geboortedatum = scanner.Text()

		//geslacht toevoegen
		fmt.Println("Wat is het geslacht van de klant? Gelieve te antwoorden met M voor man, V voor vrouw en O voor onzijdig")
		scanner.Scan()
		nieuweklant.geslacht = scanner.Text()

		//bloedgroep toevoegen
		fmt.Println("Wat is de bloedgroep van de klant?")
		scanner.Scan()
		nieuweklant.bloedgroep = scanner.Text()

		//rhesusfactor toevoegen
		fmt.Println("Wat is de rhesusfactor van de klant?")
		scanner.Scan()
		nieuweklant.rhesusfactor = scanner.Text()

		//Beroepsrisicofactor verder uitwerken.
		fmt.Println("Wat is de beroepsrisicofactor van de klant? Dit kan ingeschaald worden van 0 tot 10.")
		scanner.Scan()
		beroepsrisicofactor := scanner.Text()
		nieuweklant.beroepsrisicofactor, err = strconv.ParseFloat(beroepsrisicofactor, 64)

		//Inkomen opvragen
		fmt.Println("Wat is het inkomen van de klant?")
		scanner.Scan()
		inkomen := scanner.Text()
		nieuweklant.inkomen, err = strconv.Atoi(inkomen)

		//Kredietregistratie toevoegen
		fmt.Println("Heeft de klant een kredietregistratie? Gelieve te antwoorden met een J voor Ja en een N voor Nee")
		scanner.Scan()
		nieuweklant.kredietregistratie = scanner.Text()

		//Opleiding invoeren
		fmt.Println("Wat voor opleiding heeft de klant genoten? Kies uit WO, HBO, VMBO, NVT")
		scanner.Scan()
		nieuweklant.opleiding = scanner.Text()

		//Opmerking toevoegen
		fmt.Println("Zijn er nog overige bijzonderheden?")
		scanner.Scan()
		nieuweklant.opmerking = scanner.Text()

		fmt.Println("")

		//klantnummer aanmaken
		request_klantnummer, err := db.Query("SELECT MAX(klantnummer) FROM vitaintellectdb.klant")
		if err != nil {
			panic(err.Error())
		}

		var hoogsteklantnummer Klant
		if request_klantnummer.Next() {
			err := request_klantnummer.Scan(&hoogsteklantnummer.klantnummer)
			if err != nil {
				panic(err)
			}
		}

		nieuweklant.klantnummer = hoogsteklantnummer.klantnummer + 1

		request_bestelnummer, err := db.Query("SELECT MAX(bestelnummer) FROM vitaintellectdb.bestelling")
		if err != nil {
			panic(err.Error())
		}

		var hoogstebestelnummer bestellingen_s
		if request_bestelnummer.Next() {
			err := request_bestelnummer.Scan(&hoogstebestelnummer.bestelnummer)
			if err != nil {
				panic(err)
			}
		}

		//oude hoogste bestelnummer opgevraagd en hier 1 bij optellen
		nieuwe_bestelnummer := hoogstebestelnummer.bestelnummer + 1

		insert, err := db.Query("INSERT INTO vitaintellectdb.klant (klantnummer, naam, voornaam, postcode, huisnummer, huisnummer_toevoeging, geboortedatum, geslacht, bloedgroep, rhesusfactor, beroepsrisicofactor, inkomen, kredietregistratie, opleiding, opmerkingen)VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);", nieuweklant.klantnummer, nieuweklant.naam, nieuweklant.voornaam, nieuweklant.postcode, nieuweklant.huisnummer, nieuweklant.huisnummer_toevoeging, nieuweklant.geboortedatum, nieuweklant.geslacht, nieuweklant.bloedgroep, nieuweklant.rhesusfactor, nieuweklant.beroepsrisicofactor, nieuweklant.inkomen, nieuweklant.kredietregistratie, nieuweklant.opleiding, nieuweklant.opmerking)

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()

		//modules toevoegen en de prijs berekenen
		for {
			modulesvandeklant := modulestoevoegen()

			var totaalprijs float64
			for _, f := range modulesvandeklant {
				totaalprijs += f.stukprijs
			}

			//looptijd bepalen doormiddel van functie bepalenlooptijd()
			looptijd := bepalenlooptijd(nieuweklant.klantnummer, nieuweklant.beroepsrisicofactor)

			//hier word het maandbedrag berekend zonder kredietregistratie of inkomen
			maandbedrag := totaalprijs / looptijd

			maandbedrag_float_afgerond := fmt.Sprintf("%.2f", maandbedrag)
			maandbedraggoed, _ := strconv.ParseFloat(maandbedrag_float_afgerond, 2)

			maximalemaandbedrag := maxmaandbedrag(nieuweklant.inkomen, nieuweklant.kredietregistratie)

			if maandbedrag > maximalemaandbedrag {
				fmt.Println("Het is niet mogelijk om deze configuratie aan te schaffen gezien het inkomen en andere factoren")
			}

			fmt.Printf("Zie onderstaand de bestelling: \n\n")
			fmt.Printf(" ------------------------- ------------------------- \n")
			fmt.Printf("|%-25s|%-25s|\n", "Naam: ", nieuweklant.naam)
			fmt.Printf(" ------------------------- ------------------------- \n")
			fmt.Printf("|%-25s|%-25s|\n", "Module", "Prijs")
			fmt.Printf(" ------------------------- ------------------------- \n")

			for _, f := range modulesvandeklant {
				fmt.Printf("|%-25s|%-25g|\n", f.modulenaam, f.stukprijs)
			}
			fmt.Printf(" ------------------------- ------------------------- \n")
			fmt.Printf("|%-25s|%-25g|\n", "Totaal", totaalprijs)
			fmt.Printf(" ------------------------- ------------------------- \n")
			fmt.Printf("|%-25s|%-25g|\n", "Maximale maandtermijn", maximalemaandbedrag)
			fmt.Printf("|%-25s|%-25g|\n", "Maandtermijn", maandbedraggoed)
			fmt.Printf("|%-25s|%-25g|\n", "Looptijd in maanden", looptijd)
			fmt.Printf(" ------------------------- ------------------------- \n")

			fmt.Println("Wil je deze bestelling bevestigen? Typ Ja voor bevestiging en Nee om opnieuw te beginnen")
			scanner.Scan()
			antwoord := scanner.Text()

			current_date := time.Now()
			doorlooptijdint := int64(looptijd) + 1

			query := "SELECT naam, datum_in_dienst, medewerkernummer FROM vitaintellectdb.medewerker WHERE naam = ? and datum_in_dienst = ? and functie LIKE '%verkoper%'"

			dbresultaat, err := db.Query(query, gebruikersnaam, datum_in_dienst)
			if err != nil {
				panic(err)
			}

			var verkoper medewerker
			if dbresultaat.Next() {
				err := dbresultaat.Scan(&verkoper.naam, &verkoper.datum_in_dienst, &verkoper.medewerkernummer)
				if err != nil {
					panic(err)
				}
			}

			if antwoord == "Ja" {
				insert, err := db.Query("INSERT INTO vitaintellectdb.bestelling (bestelnummer, status, besteldatum, afbetaling_doorlooptijd, afbetaling_maandbedrag, klantnummer, verkoper)VALUES (?,?,?,?,?,?,?);", nieuwe_bestelnummer, "OFF", current_date, doorlooptijdint, maandbedraggoed, nieuweklant.klantnummer, verkoper.medewerkernummer)

				// if there is an error inserting, handle it
				if err != nil {
					panic(err.Error())
				}
				// be careful deferring Queries if you are using transactions
				defer insert.Close()
				fmt.Printf("De bestelling is opgemaakt \n\n\n")
				break loop

			} else if antwoord == "Nee" {
				fmt.Println("Wilt u het bestel proces stoppen? Typ dan Ja, wilt u opnieuw de modules toevoegen typ dan Nee.")
				scanner.Scan()
				antwoord := scanner.Text()
				if antwoord == "Ja" {
					break loop
				} else {
					fmt.Println("U kunt nu opnieuw modules toevoegen")
				}
			}

		}
	}
}

//bestaande klanten met nieuwe bestelling
func bestaandeklant(gebruikersnaam string, datum_in_dienst string) {

	scanner := bufio.NewScanner(os.Stdin)

	//Gegevens van de klant opvragen
	fmt.Println("Wat is de achternaam van de klant? ")
	scanner.Scan()
	achternaam := scanner.Text()
	fmt.Println("Wat de is de geboorte datum van de klant? YYYY-MM-DD")
	scanner.Scan()
	geboortedatum := scanner.Text()

	//hulp variabelen
	nvtfloat := 0.0
	nvtint := 1
	nvtstring := "nvt"

	db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "SELECT klantnummer, naam, IFNULL(beroepsrisicofactor, ?) as beroepsrisicofactor, IFNULL(inkomen, ?) as inkomen, IFNULL(kredietregistratie, ?) as kredietregistratie from vitaintellectdb.klant WHERE naam = ? and geboortedatum = ?"

	dbresultaat, err := db.Query(query, nvtfloat, nvtint, nvtstring, achternaam, geboortedatum)
	if err != nil {
		panic(err)
	}

	var deklant Klant
	if dbresultaat.Next() {
		err := dbresultaat.Scan(&deklant.klantnummer, &deklant.naam, &deklant.beroepsrisicofactor, &deklant.inkomen, &deklant.kredietregistratie)
		if err != nil {
			panic(err)
		}
	}

	//controleren van de beroepsrisico registratie in de database
	if deklant.beroepsrisicofactor == nvtfloat {
		fmt.Printf("Op dit moment is er nog niks over de beroepsrisicofactor bekend van %s \n", deklant.naam)
		fmt.Println("Gelieve de beroepsrisicofactor toe te voegen: ")
		scanner.Scan()
		beroepsrisicofactor := scanner.Text()
		nieuweberoepsrisicofactor, _ := strconv.ParseFloat(beroepsrisicofactor, 64)
		deklant.beroepsrisicofactor = nieuweberoepsrisicofactor

		insert, err := db.Query("UPDATE vitaintellectdb.klant SET beroepsrisicofactor = ? WHERE klantnummer = ?", deklant.beroepsrisicofactor, deklant.klantnummer)

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()
		fmt.Println("De beroepsrisicofactor is geupdate")
	}

	//controleen van het inkomen van de klant in de database
	if deklant.inkomen == nvtint {
		fmt.Printf("Op dit moment is er nog niks het inkomen van %s \n", deklant.naam)
		fmt.Println("Gelieve het inkomen toe te voegen: ")
		scanner.Scan()
		inkomen := scanner.Text()
		deklant.inkomen, err = strconv.Atoi(inkomen)
		insert, err := db.Query("UPDATE vitaintellectdb.klant SET inkomen = ? WHERE klantnummer = ?", deklant.inkomen, deklant.klantnummer)

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()
		fmt.Println("Het inkomen is geupdate")
	}

	//controleren van de kredietregistratie in de datbase
	if deklant.kredietregistratie == nvtstring {
		fmt.Printf("Op dit moment is er nog niks over de kredietregistratie bekend van %s \n", deklant.naam)
		fmt.Println("Gelieve aan te geven of de klant een kredietregistratie heeft. Typ J voor Ja en typ N voor nee")
		scanner.Scan()
		deklant.kredietregistratie = scanner.Text()
		insert, err := db.Query("UPDATE vitaintellectdb.klant SET kredietregistratie = ? WHERE klantnummer = ?", deklant.kredietregistratie, deklant.klantnummer)

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()
		fmt.Println("De kredietregistratie is toegevoegd")
	}

	//query voor het opmaken van het bestelnummer
	request_bestelnummer, err := db.Query("SELECT MAX(bestelnummer) FROM vitaintellectdb.bestelling")
	if err != nil {
		panic(err.Error())
	}

	var hoogstebestelnummer bestellingen_s
	if request_bestelnummer.Next() {
		err := request_bestelnummer.Scan(&hoogstebestelnummer.bestelnummer)
		if err != nil {
			panic(err)
		}
	}
	//nieuw bestelnummer creeëren
	nieuwe_bestelnummer := hoogstebestelnummer.bestelnummer + 1

	fmt.Println("We gaan nu een nieuwproduct toevoegen voor: ", achternaam)
	fmt.Println("")

	//toevoegen van de modules
	for {

		//starten van de modules toevoegen functie
		modulesvandeklant := modulestoevoegen()

		//hulp variabel
		var totaalprijs float64

		//de totale prijs berekenen
		for _, f := range modulesvandeklant {
			totaalprijs += f.stukprijs
		}

		//looptijd berekenen met behulp van de bepalenleeftijd functie
		looptijd := bepalenlooptijd(deklant.klantnummer, deklant.beroepsrisicofactor)

		//maandbedrag berekenen
		maandbedrag := totaalprijs / looptijd

		maandbedrag_float_afgerond := fmt.Sprintf("%.2f", maandbedrag)
		maandbedraggoed, _ := strconv.ParseFloat(maandbedrag_float_afgerond, 2)

		//maximalemaandbedrag berekenen
		maximalemaandbedrag := maxmaandbedrag(deklant.inkomen, deklant.kredietregistratie)

		fmt.Println("")

		//wanneer het maandbedrag hoger is dan het maximale berekende maandbedrag dan kan de bestelling niet voortgezet worden
		if maandbedrag > maximalemaandbedrag {
			fmt.Printf("Het is niet mogelijk om deze configuratie aan te schaffen gezien het inkomen en andere factoren\n\n")

		}

		//overzicht creeeren bestelling van de klant
		fmt.Printf("Zie onderstaand de bestelling: \n\n")
		fmt.Printf(" ------------------------- ------------------------- \n")
		fmt.Printf("|%-25s|%-25s|\n", "Naam: ", deklant.naam)
		fmt.Printf(" ------------------------- ------------------------- \n")
		fmt.Printf("|%-25s|%-25s|\n", "Module", "Prijs")
		fmt.Printf(" ------------------------- ------------------------- \n")

		for _, f := range modulesvandeklant {
			fmt.Printf("|%-25s|%-25g|\n", f.modulenaam, f.stukprijs)
		}
		fmt.Printf(" ------------------------- ------------------------- \n")
		fmt.Printf("|%-25s|%-25g|\n", "Totaal", totaalprijs)
		fmt.Printf(" ------------------------- ------------------------- \n")
		fmt.Printf("|%-25s|%-25g|\n", "Maximale maandtermijn", maximalemaandbedrag)
		fmt.Printf("|%-25s|%-25g|\n", "Maandtermijn", maandbedraggoed)
		fmt.Printf("|%-25s|%-25g|\n", "Looptijd in maanden", looptijd)
		fmt.Printf(" ------------------------- ------------------------- \n")

		//bevestigingsvraag voor de verkoper
		fmt.Println("Wil je deze bestelling bevestigen? Typ Ja voor bevestiging en Nee om opnieuw te beginnen")
		scanner.Scan()
		antwoord := scanner.Text()

		//opvragen van de huidige tijd om deze te kunnen gebruiken bij het plaatsen van de order
		current_date := time.Now()

		//doorlooptijd bepalen in int, de +1 is omdat altijd naar beneden wordt afgerond
		doorlooptijdint := int64(looptijd) + 1

		//de verkoper zijn gegevens laden
		query := "SELECT naam, datum_in_dienst, medewerkernummer FROM vitaintellectdb.medewerker WHERE naam = ? and datum_in_dienst = ? and functie LIKE '%verkoper%'"

		dbresultaat, err := db.Query(query, gebruikersnaam, datum_in_dienst)
		if err != nil {
			panic(err)
		}

		var verkoper medewerker
		if dbresultaat.Next() {
			err := dbresultaat.Scan(&verkoper.naam, &verkoper.datum_in_dienst, &verkoper.medewerkernummer)
			if err != nil {
				panic(err)
			}
		}

		//bestelling is database plaatsen
		if antwoord == "Ja" {
			insert, err := db.Query("INSERT INTO vitaintellectdb.bestelling (bestelnummer, status, besteldatum, afbetaling_doorlooptijd, afbetaling_maandbedrag, klantnummer, verkoper)VALUES (?,?,?,?,?,?,?);", nieuwe_bestelnummer, "OFF", current_date, doorlooptijdint, maandbedraggoed, deklant.klantnummer, verkoper.medewerkernummer)

			// if there is an error inserting, handle it
			if err != nil {
				panic(err.Error())
			}
			// be careful deferring Queries if you are using transactions
			defer insert.Close()
			fmt.Printf("De bestelling is opgemaakt \n\n\n")
			break

		} else if antwoord == "Nee" {
			fmt.Println("Wilt u het bestel proces stoppen? Typ dan Ja, wilt u opnieuw de modules toevoegen typ dan Nee.")
			scanner.Scan()
			antwoord := scanner.Text()
			if antwoord == "Ja" {
				break
			} else {
				fmt.Println("U kunt nu opnieuw modules toevoegen")
			}
		}

	}
}

func modulestoevoegen() []module {

	//twee lijsten van structs
	modulenaam := []module{}
	modulekeuzeklant := []module{}

	//database connectie opzetten
	db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")
	if err != nil {
		panic(err.Error())
	}

	//opvragen van de modules die verkocht worden door Vita Intelectus in de database
	results, err := db.Query("SELECT modulenaam, stukprijs from vitaintellectdb.module")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var modulenaam1 module
		// for each row, scan the result into our tag composite object
		err = results.Scan(&modulenaam1.modulenaam, &modulenaam1.stukprijs)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		modulenaam = append(modulenaam, modulenaam1)
		if modulenaam1.modulenaam == "basis" {
			modulekeuzeklant = append(modulekeuzeklant, modulenaam1)
		}
	}

	fmt.Println("De klant heeft altijd een basismodule nodig om vervolgens uit te breiden met andere modules. Daarom is standaard de basismodule al geselecteerd.")

	fmt.Println("")

	//het toevoegen van de modules
Loop:
	for {
		fmt.Println(" ")
		fmt.Println("Welke module wil je graag toevoegen voor de klant?")
		fmt.Println(" ")
		fmt.Println("Toets 1 voor de Cor uitbreiding")
		fmt.Println("Toets 2 voor de Dermal uitbreiding")
		fmt.Println("Toets 3 voor de Memoria uitbreiding")
		fmt.Println("Toets 4 voor de Oculus uitbreiding")
		fmt.Println("Toets 5 voor de Oricula uitbreiding")
		fmt.Println("Toets 6 voor de Pes uitbreiding")
		fmt.Println("Toets 7 voor de Sangius uitbreiding")
		fmt.Println("Toets 8 voor de Somnus uitbreiding")
		fmt.Println("Toets 9 om verder te gaan naar het overzicht van bestelling")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		keuzevanmedewerker := scanner.Text()

		//keuze menu van de modules, wanneer keuze is gemaakt wordt die toegevoegd aan de list.
		switch keuzevanmedewerker {
		case "1":
			for _, f := range modulenaam {
				if f.modulenaam == "Cor" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Cor is toegevoegd")
		case "2":
			for _, f := range modulenaam {
				if f.modulenaam == "Dermal" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Dermal is toegevoegd")
		case "3":
			for _, f := range modulenaam {
				if f.modulenaam == "Memoria" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Memoria is toegevoegd")
		case "4":
			for _, f := range modulenaam {
				if f.modulenaam == "Oculus" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Oculus is toegevoegd")
		case "5":
			for _, f := range modulenaam {
				if f.modulenaam == "Oricula" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Oricula is toegevoegd")
		case "6":
			for _, f := range modulenaam {
				if f.modulenaam == "Pes" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Pes is toegevoegd")
		case "7":
			for _, f := range modulenaam {
				if f.modulenaam == "Sangius" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Sangius is toegevoegd")
		case "8":
			for _, f := range modulenaam {
				if f.modulenaam == "Somnius" {
					modulekeuzeklant = append(modulekeuzeklant, f)
				}
			}
			fmt.Println("De module Somnius is toegevoegd")
		case "9":
			break Loop
		}

		fmt.Printf("De volgende modules zijn voor de klant geselecteerd: ")
		for _, f := range modulekeuzeklant {
			fmt.Printf("%v ", f.modulenaam)
		}
	}

	return modulekeuzeklant

}

func maxmaandbedrag(inkomen int, krediet string) float64 {

	//Voor de opdracht ben ik er vanuit gegaan dat bij een kredietregistratie de klant 80% van het normale maandbedrag kan lenen.
	hetinkomen := inkomen
	inkomenfloat := float64(hetinkomen)
	deler := 120
	delerfloat := float64(deler)

	var maximaalmaandbedrag float64
	if krediet == "N" {
		maximaalmaandbedrag = inkomenfloat / delerfloat
	} else if krediet == "Y" {
		maximaalmaandbedrag = inkomenfloat / delerfloat * 0.8
	}

	return maximaalmaandbedrag

}

func bepalenleeftijd(klantnummer int) int {
	//hier word de leeftijd van de klant opgehaald uit de database en gereturned
	//struct van leeftijd
	type age struct {
		leeftijd int
	}

	db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//query opstellen
	query := "SELECT timestampdiff(YEAR, geboortedatum, CURDATE()) FROM vitaintellectdb.klant as AGE WHERE klantnummer=?;"
	resultaat, err := db.Query(query, klantnummer)
	if err != nil {
		panic(err)
	}

	var deleeftijd age
	if resultaat.Next() {
		err := resultaat.Scan(&deleeftijd.leeftijd)
		if err != nil {
			panic(err)
		}
	}
	return deleeftijd.leeftijd
}

//Looptijd bepalen
func bepalenlooptijd(klantnummer int, beroepsrisico float64) float64 {
	//hier wordt de looptijd berekend door middel van de leeftijd en beroeprisico factor

	//leeftijd wordt bepaalt doormiddel van voorop gestelde functie
	leeftijd := bepalenleeftijd(klantnummer)

	//bepalen van beroepsrisico maanden
	beroepsrisicofactor := beroepsrisico / 1.5 * 12

	var looptijd_in_maand_leeftijd float64
	if leeftijd < 45 {
		looptijd_in_maand_leeftijd = 120
	} else if leeftijd >= 45 && leeftijd < 55 {
		looptijd_in_maand_leeftijd = 90
	} else if leeftijd >= 55 {
		looptijd_in_maand_leeftijd = 60
	}
	//berekenen van de looptijd rekeninghoudend met het beroepsrisico
	looptijd := looptijd_in_maand_leeftijd - beroepsrisicofactor

	return looptijd

}

//Geeft alle klanten weer
func AlleKlanten() {

	//Het is mogelijk om alles kolommen uit de database weer te geven, dit maakt het echter denk ik onoverzichtelijk.
	//Dan moet er met een SQL query alle data opgehaald worden en de functie IFNUL() gebruiken
	db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")
	if err != nil {
		panic(err.Error())
	}

	results, err := db.Query("SELECT klantnummer, naam, voornaam FROM vitaintellectdb.klant")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Printf(" --------------- --------------- --------------- ")
	fmt.Println("")
	fmt.Printf("|%-15s|%-15s|%-15s|\n", "Klantnummer ", "Achternaam ", "Voornaam ")
	fmt.Printf(" --------------- --------------- --------------- ")
	fmt.Println("")

	for results.Next() {
		var klant Klant
		// for each row, scan the result into our tag composite object
		err = results.Scan(&klant.klantnummer, &klant.naam, &klant.voornaam)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		//fmt.Println(klant.klantnummer, "		" , klant.voornaam, "			" , klant.naam)
		fmt.Printf("|%-15d|%-15s|%-15s|\n", klant.klantnummer, klant.naam, klant.voornaam)
	}

	fmt.Printf(" --------------- --------------- --------------- ")
	fmt.Println("")
	defer db.Close()
}

//Geeft alle openstaande order weer
func AllOpenOrders() {

	//Dan moet er met een SQL query alle data opgehaald worden en de functie IFNUL() gebruiken
	fmt.Printf(" --------------- --------------- --------------- --------------- ----------------")
	fmt.Println("")
	fmt.Printf("|%-15s|%-15s|%-15s|%-15s|%-15s|\n", "Bestelnumnmer", "Bestel Status", "Besteldatum", "klantnummer", "Medewerkernummer")
	fmt.Printf(" --------------- --------------- --------------- --------------- ----------------")
	fmt.Println("")

	db, err := sql.Open("mysql", "root:Wachtwoord1!@tcp(127.0.0.1:3306)/vitaintellectdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM vitaintellectdb.bestelling")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var bestellingen_s bestellingen_s
		// for each row, scan the result into our tag composite object
		err = results.Scan(&bestellingen_s.bestelnummer, &bestellingen_s.status, &bestellingen_s.besteldatum, &bestellingen_s.afbetaling_doorlooptijd, &bestellingen_s.afbetaling_maandbedrag, &bestellingen_s.klantnummer, &bestellingen_s.verkoper)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		fmt.Printf("|%-15d|%-15s|%-15s|%-15d|%-16d|\n", bestellingen_s.bestelnummer, bestellingen_s.status, bestellingen_s.besteldatum, bestellingen_s.klantnummer, bestellingen_s.verkoper)
	}
	fmt.Printf(" --------------- --------------- --------------- --------------- ----------------")
	fmt.Println("")
	defer db.Close()

}
