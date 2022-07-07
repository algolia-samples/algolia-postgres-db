package index

import (
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

var (
	algolia_app_id = os.Getenv("ALGOLIA_APP_ID")
	algolia_index  = os.Getenv("ALGOLIA_INDEX_NAME")
	algolia_apikey = os.Getenv("ALGOLIA_API_KEY")
)

func UploadRecordsToAlgola(records []AuditLogRecord) error {
	// Connect and authenticate with your Algolia app
	client := search.NewClient(algolia_app_id, algolia_apikey)
	// Create a new index and add a record
	index := client.InitIndex(algolia_index)
	resSave, err := index.SaveObjects(records)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// ids := "Processed IDS: "
	// for _, record := range records {
	// 	ids += fmt.Sprintf("%d, ", record.Id)
	// }
	// log.Println(ids)
	resSave.Wait()

	// Search the index and print the results
	indexContent, err := index.Search("")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Found %d records in the index.", indexContent.NbHits)

	return err
}
