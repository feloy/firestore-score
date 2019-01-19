# Save scores into Firestore, using Go Functions

More information on https://medium.com/@feloy/firebase-saving-scores-into-firestore-with-go-functions-b128fd8c425

To publish the function:

```
gcloud beta functions deploy newScore \
    --entry-point OnNewScore \
    --memory 256MB \
    --region us-central1 \
    --runtime go111 \
    --trigger-event providers/cloud.firestore/eventTypes/document.create \
    --trigger-resource "projects/<your-project-id>/databases/(default)/documents/Score/{id}"
```
