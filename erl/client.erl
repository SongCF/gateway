%%%-------------------------------------------------------------------
%%% @author 'SongCF'
%%% @copyright (C) 2016, <JH>
%%% @doc
%%%
%%% @end
%%% Created : 24. 十月 2016 16:28
%%%-------------------------------------------------------------------
-module(client).

%% INCLUDE

%% EXPORT
-compile([export_all]).


-behaviour(gen_server).

-record(state, {sock, pid}).
-define(HEARTBEAT_TIME, 120000). %2m


%% ------------------------------------------------------------------
%% API Function Definitions
%% ------------------------------------------------------------------

%% IP::string()
%% Port::integer()
start(IP, Port) ->
    {ok, Pid} = gen_server:start_link(?MODULE, [IP, Port], []),
    Pid.

stop(Pid) ->
    Pid ! stop,
    timer:sleep(1000).

call(Pid, Msg)->
    gen_server:call(Pid, Msg).


%% ------------------------------------------------------------------
%% gen_server Function Definitions
%% ------------------------------------------------------------------

init([IP, Port]) ->
    {ok, Socket} = gen_tcp:connect(IP, Port, [{packet, 2}, {active, true}]),
    erlang:send_after(?HEARTBEAT_TIME, self(), heartbeat),
    {ok, #state{sock = Socket}}.

handle_call({send, Data}, _From, #state{sock = Socket}=State) ->
    io:format("send ..."),
    timer:sleep(2000),
    gen_tcp:send(Socket, Data),
    {reply, ok, State};

handle_call(_Request, _From, State) ->
    {reply, ok, State}.

handle_cast(_Msg, State) ->
    {noreply, State}.

handle_info({tcp, _Socket, Msg}, #state{pid=Pid}=State) ->
    io:format("Message = ~p~n", [Msg]),
    {noreply, State};

handle_info(heartbeat, #state{sock = Socket}=State) ->
    gen_tcp:send(Socket, "heartbeat"),
    erlang:send_after(?HEARTBEAT_TIME, self(), heartbeat),
    {noreply, State};

handle_info(stop, State) ->
    {stop, normal, State};

handle_info(Info, State) ->
    io:format("unknown Info = ~p~n", [Info]),
    {noreply, State}.

terminate(Reason, _State) ->
    io:format("terminate ~p [~p]", [?MODULE, Reason]),
    ok.

code_change(_OldVsn, State, _Extra) ->
    {ok, State}.

