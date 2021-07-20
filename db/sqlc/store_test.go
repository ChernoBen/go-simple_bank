package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//usando goroutine
	n := 5
	amount := int64(10)
	//usando chanels para verificar retorno da operação
	errs := make(chan error)               //channel para receber errors
	results := make(chan TransferTxResult) //channel para receber resultados

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TranferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			//enviando o erro para o channel de errors e resultados para channel de resultados
			errs <- err
			results <- result

		}()
	}
	//verificando retornos recebidos pelos channels
	for i := 0; i < n; i++ {
		//obtendo erro de dentro do channel
		err := <-errs
		require.NoError(t, err) //verificar se existe erro dentro do channel errs
		//obtendo resultados de dentro do channel resultados
		result := <-results
		require.NotEmpty(t, result) //verifica se o retorno do channel nao é vazio
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntries(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntries(context.Background(), toEntry.ID)
		require.NoError(t, err)

	}
}
