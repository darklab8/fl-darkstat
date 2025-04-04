# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import darkstat_pb2 as darkstat__pb2

GRPC_GENERATED_VERSION = '1.70.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in darkstat_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class DarkstatStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetHealth = channel.unary_unary(
                '/statproto.Darkstat/GetHealth',
                request_serializer=darkstat__pb2.Empty.SerializeToString,
                response_deserializer=darkstat__pb2.HealthReply.FromString,
                _registered_method=True)
        self.GetBasesNpc = channel.unary_unary(
                '/statproto.Darkstat/GetBasesNpc',
                request_serializer=darkstat__pb2.GetBasesInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetBasesReply.FromString,
                _registered_method=True)
        self.GetBasesMiningOperations = channel.unary_unary(
                '/statproto.Darkstat/GetBasesMiningOperations',
                request_serializer=darkstat__pb2.GetBasesInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetBasesReply.FromString,
                _registered_method=True)
        self.GetBasesPoBs = channel.unary_unary(
                '/statproto.Darkstat/GetBasesPoBs',
                request_serializer=darkstat__pb2.GetBasesInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetBasesReply.FromString,
                _registered_method=True)
        self.GetPoBs = channel.unary_unary(
                '/statproto.Darkstat/GetPoBs',
                request_serializer=darkstat__pb2.Empty.SerializeToString,
                response_deserializer=darkstat__pb2.GetPoBsReply.FromString,
                _registered_method=True)
        self.GetPoBGoods = channel.unary_unary(
                '/statproto.Darkstat/GetPoBGoods',
                request_serializer=darkstat__pb2.Empty.SerializeToString,
                response_deserializer=darkstat__pb2.GetPoBGoodsReply.FromString,
                _registered_method=True)
        self.GetCommodities = channel.unary_unary(
                '/statproto.Darkstat/GetCommodities',
                request_serializer=darkstat__pb2.GetCommoditiesInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetCommoditiesReply.FromString,
                _registered_method=True)
        self.GetGuns = channel.unary_unary(
                '/statproto.Darkstat/GetGuns',
                request_serializer=darkstat__pb2.GetGunsInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetGunsReply.FromString,
                _registered_method=True)
        self.GetMissiles = channel.unary_unary(
                '/statproto.Darkstat/GetMissiles',
                request_serializer=darkstat__pb2.GetGunsInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetGunsReply.FromString,
                _registered_method=True)
        self.GetAmmos = channel.unary_unary(
                '/statproto.Darkstat/GetAmmos',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetAmmoReply.FromString,
                _registered_method=True)
        self.GetCounterMeasures = channel.unary_unary(
                '/statproto.Darkstat/GetCounterMeasures',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetCounterMeasuresReply.FromString,
                _registered_method=True)
        self.GetEngines = channel.unary_unary(
                '/statproto.Darkstat/GetEngines',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetEnginesReply.FromString,
                _registered_method=True)
        self.GetMines = channel.unary_unary(
                '/statproto.Darkstat/GetMines',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetMinesReply.FromString,
                _registered_method=True)
        self.GetScanners = channel.unary_unary(
                '/statproto.Darkstat/GetScanners',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetScannersReply.FromString,
                _registered_method=True)
        self.GetShields = channel.unary_unary(
                '/statproto.Darkstat/GetShields',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetShieldsReply.FromString,
                _registered_method=True)
        self.GetShips = channel.unary_unary(
                '/statproto.Darkstat/GetShips',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetShipsReply.FromString,
                _registered_method=True)
        self.GetThrusters = channel.unary_unary(
                '/statproto.Darkstat/GetThrusters',
                request_serializer=darkstat__pb2.GetEquipmentInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetThrustersReply.FromString,
                _registered_method=True)
        self.GetFactions = channel.unary_unary(
                '/statproto.Darkstat/GetFactions',
                request_serializer=darkstat__pb2.GetFactionsInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetFactionsReply.FromString,
                _registered_method=True)
        self.GetTractors = channel.unary_unary(
                '/statproto.Darkstat/GetTractors',
                request_serializer=darkstat__pb2.GetTractorsInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetTractorsReply.FromString,
                _registered_method=True)
        self.GetHashes = channel.unary_unary(
                '/statproto.Darkstat/GetHashes',
                request_serializer=darkstat__pb2.Empty.SerializeToString,
                response_deserializer=darkstat__pb2.GetHashesReply.FromString,
                _registered_method=True)
        self.GetInfocards = channel.unary_unary(
                '/statproto.Darkstat/GetInfocards',
                request_serializer=darkstat__pb2.GetInfocardsInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetInfocardsReply.FromString,
                _registered_method=True)
        self.GetGraphPaths = channel.unary_unary(
                '/statproto.Darkstat/GetGraphPaths',
                request_serializer=darkstat__pb2.GetGraphPathsInput.SerializeToString,
                response_deserializer=darkstat__pb2.GetGraphPathsReply.FromString,
                _registered_method=True)


class DarkstatServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetHealth(self, request, context):
        """Just to check if grpc works. Returns boolean value if it is healthy as true
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetBasesNpc(self, request, context):
        """Get all Freelancer NPC bases
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetBasesMiningOperations(self, request, context):
        """Get all imaginary bases that in place of mining fields. Useful for trading calculations
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetBasesPoBs(self, request, context):
        """Get all Player Owned bases in the same format as Npc bases. Returns only PoBs which have known positions
        Useful for trading calculations
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetPoBs(self, request, context):
        """Get all Player Owned bases. Completely all that are public exposed
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetPoBGoods(self, request, context):
        """Get all PoB goods, where they are sold and bought. Reverse search by PoBs
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetCommodities(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetGuns(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetMissiles(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetAmmos(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetCounterMeasures(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetEngines(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetMines(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetScanners(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetShields(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetShips(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetThrusters(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetFactions(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetTractors(self, request, context):
        """Get Tractors. For Discovery those are IDs
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetHashes(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetInfocards(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetGraphPaths(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_DarkstatServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetHealth': grpc.unary_unary_rpc_method_handler(
                    servicer.GetHealth,
                    request_deserializer=darkstat__pb2.Empty.FromString,
                    response_serializer=darkstat__pb2.HealthReply.SerializeToString,
            ),
            'GetBasesNpc': grpc.unary_unary_rpc_method_handler(
                    servicer.GetBasesNpc,
                    request_deserializer=darkstat__pb2.GetBasesInput.FromString,
                    response_serializer=darkstat__pb2.GetBasesReply.SerializeToString,
            ),
            'GetBasesMiningOperations': grpc.unary_unary_rpc_method_handler(
                    servicer.GetBasesMiningOperations,
                    request_deserializer=darkstat__pb2.GetBasesInput.FromString,
                    response_serializer=darkstat__pb2.GetBasesReply.SerializeToString,
            ),
            'GetBasesPoBs': grpc.unary_unary_rpc_method_handler(
                    servicer.GetBasesPoBs,
                    request_deserializer=darkstat__pb2.GetBasesInput.FromString,
                    response_serializer=darkstat__pb2.GetBasesReply.SerializeToString,
            ),
            'GetPoBs': grpc.unary_unary_rpc_method_handler(
                    servicer.GetPoBs,
                    request_deserializer=darkstat__pb2.Empty.FromString,
                    response_serializer=darkstat__pb2.GetPoBsReply.SerializeToString,
            ),
            'GetPoBGoods': grpc.unary_unary_rpc_method_handler(
                    servicer.GetPoBGoods,
                    request_deserializer=darkstat__pb2.Empty.FromString,
                    response_serializer=darkstat__pb2.GetPoBGoodsReply.SerializeToString,
            ),
            'GetCommodities': grpc.unary_unary_rpc_method_handler(
                    servicer.GetCommodities,
                    request_deserializer=darkstat__pb2.GetCommoditiesInput.FromString,
                    response_serializer=darkstat__pb2.GetCommoditiesReply.SerializeToString,
            ),
            'GetGuns': grpc.unary_unary_rpc_method_handler(
                    servicer.GetGuns,
                    request_deserializer=darkstat__pb2.GetGunsInput.FromString,
                    response_serializer=darkstat__pb2.GetGunsReply.SerializeToString,
            ),
            'GetMissiles': grpc.unary_unary_rpc_method_handler(
                    servicer.GetMissiles,
                    request_deserializer=darkstat__pb2.GetGunsInput.FromString,
                    response_serializer=darkstat__pb2.GetGunsReply.SerializeToString,
            ),
            'GetAmmos': grpc.unary_unary_rpc_method_handler(
                    servicer.GetAmmos,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetAmmoReply.SerializeToString,
            ),
            'GetCounterMeasures': grpc.unary_unary_rpc_method_handler(
                    servicer.GetCounterMeasures,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetCounterMeasuresReply.SerializeToString,
            ),
            'GetEngines': grpc.unary_unary_rpc_method_handler(
                    servicer.GetEngines,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetEnginesReply.SerializeToString,
            ),
            'GetMines': grpc.unary_unary_rpc_method_handler(
                    servicer.GetMines,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetMinesReply.SerializeToString,
            ),
            'GetScanners': grpc.unary_unary_rpc_method_handler(
                    servicer.GetScanners,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetScannersReply.SerializeToString,
            ),
            'GetShields': grpc.unary_unary_rpc_method_handler(
                    servicer.GetShields,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetShieldsReply.SerializeToString,
            ),
            'GetShips': grpc.unary_unary_rpc_method_handler(
                    servicer.GetShips,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetShipsReply.SerializeToString,
            ),
            'GetThrusters': grpc.unary_unary_rpc_method_handler(
                    servicer.GetThrusters,
                    request_deserializer=darkstat__pb2.GetEquipmentInput.FromString,
                    response_serializer=darkstat__pb2.GetThrustersReply.SerializeToString,
            ),
            'GetFactions': grpc.unary_unary_rpc_method_handler(
                    servicer.GetFactions,
                    request_deserializer=darkstat__pb2.GetFactionsInput.FromString,
                    response_serializer=darkstat__pb2.GetFactionsReply.SerializeToString,
            ),
            'GetTractors': grpc.unary_unary_rpc_method_handler(
                    servicer.GetTractors,
                    request_deserializer=darkstat__pb2.GetTractorsInput.FromString,
                    response_serializer=darkstat__pb2.GetTractorsReply.SerializeToString,
            ),
            'GetHashes': grpc.unary_unary_rpc_method_handler(
                    servicer.GetHashes,
                    request_deserializer=darkstat__pb2.Empty.FromString,
                    response_serializer=darkstat__pb2.GetHashesReply.SerializeToString,
            ),
            'GetInfocards': grpc.unary_unary_rpc_method_handler(
                    servicer.GetInfocards,
                    request_deserializer=darkstat__pb2.GetInfocardsInput.FromString,
                    response_serializer=darkstat__pb2.GetInfocardsReply.SerializeToString,
            ),
            'GetGraphPaths': grpc.unary_unary_rpc_method_handler(
                    servicer.GetGraphPaths,
                    request_deserializer=darkstat__pb2.GetGraphPathsInput.FromString,
                    response_serializer=darkstat__pb2.GetGraphPathsReply.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'statproto.Darkstat', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('statproto.Darkstat', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class Darkstat(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetHealth(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetHealth',
            darkstat__pb2.Empty.SerializeToString,
            darkstat__pb2.HealthReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetBasesNpc(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetBasesNpc',
            darkstat__pb2.GetBasesInput.SerializeToString,
            darkstat__pb2.GetBasesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetBasesMiningOperations(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetBasesMiningOperations',
            darkstat__pb2.GetBasesInput.SerializeToString,
            darkstat__pb2.GetBasesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetBasesPoBs(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetBasesPoBs',
            darkstat__pb2.GetBasesInput.SerializeToString,
            darkstat__pb2.GetBasesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetPoBs(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetPoBs',
            darkstat__pb2.Empty.SerializeToString,
            darkstat__pb2.GetPoBsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetPoBGoods(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetPoBGoods',
            darkstat__pb2.Empty.SerializeToString,
            darkstat__pb2.GetPoBGoodsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetCommodities(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetCommodities',
            darkstat__pb2.GetCommoditiesInput.SerializeToString,
            darkstat__pb2.GetCommoditiesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetGuns(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetGuns',
            darkstat__pb2.GetGunsInput.SerializeToString,
            darkstat__pb2.GetGunsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetMissiles(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetMissiles',
            darkstat__pb2.GetGunsInput.SerializeToString,
            darkstat__pb2.GetGunsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetAmmos(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetAmmos',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetAmmoReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetCounterMeasures(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetCounterMeasures',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetCounterMeasuresReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetEngines(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetEngines',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetEnginesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetMines(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetMines',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetMinesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetScanners(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetScanners',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetScannersReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetShields(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetShields',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetShieldsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetShips(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetShips',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetShipsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetThrusters(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetThrusters',
            darkstat__pb2.GetEquipmentInput.SerializeToString,
            darkstat__pb2.GetThrustersReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetFactions(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetFactions',
            darkstat__pb2.GetFactionsInput.SerializeToString,
            darkstat__pb2.GetFactionsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetTractors(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetTractors',
            darkstat__pb2.GetTractorsInput.SerializeToString,
            darkstat__pb2.GetTractorsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetHashes(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetHashes',
            darkstat__pb2.Empty.SerializeToString,
            darkstat__pb2.GetHashesReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetInfocards(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetInfocards',
            darkstat__pb2.GetInfocardsInput.SerializeToString,
            darkstat__pb2.GetInfocardsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetGraphPaths(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/statproto.Darkstat/GetGraphPaths',
            darkstat__pb2.GetGraphPathsInput.SerializeToString,
            darkstat__pb2.GetGraphPathsReply.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
