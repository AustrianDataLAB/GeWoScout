import pytest


@pytest.mark.usefixtures("cosmos_db_setup_listings")
def test_cosmosdb_fixture_data_insertion(cosmos_db_setup_listings):
    _database, container = cosmos_db_setup_listings

    # Check if data insertion was correct
    # There are 7 listings in the fixture for Vienna
    query = "SELECT * FROM c WHERE c._partitionKey = 'vienna'"
    items = list(container.query_items(query=query, enable_cross_partition_query=True))
    assert len(items) == 7

    # Should also work this way
    items = list(container.query_items(query=query, partition_key='vienna'))
    assert len(items) == 7

    query = "SELECT * FROM c WHERE c._partitionKey = 'salzburg'"
    items = list(container.query_items(query=query, partition_key='salzburg'))
    assert len(items) == 1
