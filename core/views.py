from core.models import UserRepo
from core.serializers import UserRepoSerializer
from django.http import Http404
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status


class UserRepoList(APIView):
    """
    List all userRepos, or create a new userRepos.
    """
    def get(self, request, format=None):
        userRepos = UserRepo.objects.all()
        serializer = UserRepoSerializer(userRepos, many=True)
        return Response(serializer.data)
