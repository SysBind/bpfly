module ExecSnoop exposing (main)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)

import Http
import Browser
import Time


-- MAIN

url : String
url =
    "http://localhost:7379"



main : Program () Model Msg
main =
    Browser.element
        { init = \flags -> ( { execs = [], errorMessage = Nothing }, getData )
        , view = view
        , update = update
        , subscriptions = subscriptions
        }      



-- MODEL

type alias Model =
    { execs: List String    
    ,errorMessage : Maybe String
    }



    

-- VIEW

view : Model -> Html Msg
view model =
        div [ class "jumbotron" ]
            [
             viewExecsOrError model
            ]
            
        


viewExecsOrError : Model -> Html Msg
viewExecsOrError model =
    case model.errorMessage of
        Just message ->
            viewError message

        Nothing ->
            viewExecs model.execs


viewError : String -> Html Msg
viewError errorMessage =
    let
        errorHeading =
            "Couldn't fetch execs at this time."
    in
        div []
            [ h3 [] [ text errorHeading ]
            , text ("Error: " ++ errorMessage)
            ]


viewExecs : List String -> Html Msg
viewExecs execs =
    div []
        [ h3 [] [ text "Execs" ]
        , ul [] (List.map viewExec execs)
        ]


viewExec : String -> Html Msg
viewExec exec =
    li [] [ text exec ]




-- UPDATE

type Msg
    = Tick Time.Posix 
    | DataReceived (Result Http.Error String)
   



getData : Cmd Msg
getData =
    Http.get
        { url = url ++ "/KEYS/execsnoop*"
        , expect = Http.expectString DataReceived
        }
      
update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        DataReceived (Ok dataStr) ->
            let
                data =
                    String.split "," dataStr
            in
                ( { execs = data, errorMessage = Nothing}, Cmd.none )

        DataReceived (Err _) ->
            ( model, Cmd.none )

        Tick newTime ->
            ( model, getData )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Time.every 1000 Tick


