from core.models import UserRepo, Issue
from core.serializers import UserRepoSerializer, IssueSerializer
from django.http import Http404
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework import generics
import django_filters.rest_framework

class UserRepoList(APIView):
    """
    List all userRepos, or create a new userRepos.
    """
    def get(self, request, format=None):
        userRepos = UserRepo.objects.all()
        serializer = UserRepoSerializer(userRepos, many=True)
        return Response(serializer.data)

class IssueList(generics.ListCreateAPIView):
    """
    List all Issues.
    """
    queryset = Issue.objects.all()
    serializer_class = IssueSerializer
    filter_backends = (django_filters.rest_framework.DjangoFilterBackend,)
    filter_fields = ('language', 'tech_stack', 'experience_needed', 'expected_time',)
