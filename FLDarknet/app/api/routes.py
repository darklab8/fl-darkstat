from flask import jsonify, request
from app.api import bp

@bp.route('/test', methods=['GET'])
def get_user(id):
    page = request.args.get('page', 1, type=int)
    data = request.get_json() or {}
    return jsonify('123')

@bp.route('/error', methods=['GET'])
def get_user2(id):
    page = request.args.get('page', 1, type=int)
    data = request.get_json() or {}
    return bad_request('123')


