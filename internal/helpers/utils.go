package helpers

import (
    "time"
    "math/rand"
)

const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init( ) {
    rand.Seed( time.Now( ).UnixNano( ) )
}

func getRandomString( n int )( string ) {
    rbytes := make( [ ]byte, n )

    for i, _ := range rbytes { 
        rbytes[ i ] = alphabets[ rand.Intn( len( alphabets ) ) ]
    }

    return string( rbytes )
}
